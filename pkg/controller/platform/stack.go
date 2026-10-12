package platform

import (
	"context"
	"reflect"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infinimeshv1beta1 "github.com/infinimesh/enterprise-operator/pkg/apis/infinimesh/v1beta1"
)

func coreImage(name, tag string) string {
	if tag == "" {
		tag = defaultCoreVersion
	}
	return "ghcr.io/infinimesh/core/" + name + ":" + tag
}

func imagePullPolicy(image string) corev1.PullPolicy {
	if strings.HasPrefix(image, "ghcr.io/infinimesh/core/") && !strings.HasSuffix(image, ":latest") {
		return corev1.PullIfNotPresent
	}
	return corev1.PullAlways
}

func (r *ReconcilePlatform) reconcileCore(instance *infinimeshv1beta1.Platform) error {
	if err := r.ensureAuth(instance); err != nil {
		return err
	}
	auth := instance.Name + "-auth"
	prefix := instance.Name
	secretRef := func(key string) corev1.EnvVar {
		return corev1.EnvVar{
			Name: key,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: auth},
					Key:                  key,
				},
			},
		}
	}
	work := []struct {
		suffix string
		image  string
		ports  []int
		env    []corev1.EnvVar
		cert   bool
		data   string
	}{
		{suffix: "arango", image: "arangodb/arangodb:latest", ports: []int{8529}, data: "/var/lib/arangodb3", env: []corev1.EnvVar{secretRef("ARANGO_ROOT_PASSWORD")}},
		{suffix: "rabbitmq", image: "rabbitmq:3.9-management", ports: []int{5672, 15672}, data: "/var/lib/rabbitmq", env: []corev1.EnvVar{
			{Name: "RABBITMQ_DEFAULT_USER", ValueFrom: secretRef("RABBITMQ_USER").ValueFrom},
			{Name: "RABBITMQ_DEFAULT_PASS", ValueFrom: secretRef("RABBITMQ_PASS").ValueFrom},
		}},
		{suffix: "redis", image: "redis:latest", ports: []int{6379}, data: "/data"},
		{suffix: "api", image: coreImage("api", instance.Status.Version), ports: []int{8000}, env: []corev1.EnvVar{
			{Name: "DB_HOST", Value: prefix + "-arango:8529"},
			secretRef("DB_CRED"),
			{Name: "SHADOW_HOST", Value: prefix + "-shadow:8000"},
			{Name: "PAIRING_HOST", Value: prefix + "-pairing:8000"},
			{Name: "REGISTRY_HOST", Value: prefix + "-api:8000"},
			{Name: "NAMESPACES_HOST", Value: prefix + "-api:8000"},
			{Name: "REDIS_HOST", Value: prefix + "-redis:6379"},
			secretRef("SIGNING_KEY"),
			secretRef("INF_DEFAULT_ROOT_PASS"),
		}},
		{suffix: "shadow", image: coreImage("shadow", instance.Status.Version), ports: []int{8000}, env: []corev1.EnvVar{
			secretRef("RABBITMQ_CONN"),
			{Name: "REDIS_HOST", Value: prefix + "-redis:6379"},
			{Name: "DEVICES_HOST", Value: prefix + "-api:8000"},
			secretRef("SIGNING_KEY"),
		}},
		{suffix: "mqtt-bridge", image: coreImage("mqtt-bridge", instance.Status.Version), ports: []int{8883, 2112}, cert: true, env: []corev1.EnvVar{
			secretRef("RABBITMQ_CONN"),
			{Name: "DEVICES_HOST", Value: prefix + "-api:8000"},
			{Name: "SHADOW_HOST", Value: prefix + "-shadow:8000"},
			secretRef("SIGNING_KEY"),
		}},
		{suffix: "gateway", image: coreImage("gateway", instance.Status.Version), ports: []int{8000}, env: []corev1.EnvVar{
			{Name: "APISERVER_HOST", Value: prefix + "-api:8000"},
			{Name: "PAIRING", Value: prefix + "-pairing:8000"},
		}},
		{suffix: "console", image: coreImage("console", instance.Status.Version), ports: []int{80}},
		{suffix: "pairing", image: coreImage("pairing", instance.Status.Version), ports: []int{8000}, env: []corev1.EnvVar{
			{Name: "REDIS_HOST", Value: prefix + "-redis:6379"},
			secretRef("SIGNING_KEY"),
		}},
		{suffix: "rules", image: coreImage("rules", instance.Status.Version), ports: []int{2112}, env: []corev1.EnvVar{
			{Name: "DB_HOST", Value: prefix + "-arango:8529"},
			secretRef("DB_CRED"),
			{Name: "REDIS_HOST", Value: prefix + "-redis:6379"},
			{Name: "SHADOW_HOST", Value: prefix + "-shadow:8000"},
			secretRef("INF_DEFAULT_ROOT_PASS"),
		}},
		{suffix: "coap", image: coreImage("coap-gateway", instance.Status.Version), ports: []int{5683, 2112}, env: []corev1.EnvVar{
			secretRef("RABBITMQ_CONN"),
			{Name: "COAP_ADDR", Value: ":5683"},
		}},
		{suffix: "http-fs", image: coreImage("http-fs", instance.Status.Version), ports: []int{8000, 2112}, data: "/static", env: []corev1.EnvVar{
			{Name: "DB_HOST", Value: prefix + "-arango:8529"},
			secretRef("DB_CRED"),
			secretRef("INF_DEFAULT_ROOT_PASS"),
			secretRef("SIGNING_KEY"),
			{Name: "HTTP_FS_ROOT", Value: "/static"},
		}},
	}
	for _, item := range work {
		if item.data != "" {
			if err := r.ensurePVC(instance, prefix+"-"+item.suffix); err != nil {
				return err
			}
		}
		if err := r.ensureWorkload(instance, item.suffix, item.image, item.ports, item.env, item.cert, item.data); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcilePlatform) ensureAuth(instance *infinimeshv1beta1.Platform) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: instance.Name + "-auth", Namespace: instance.Namespace},
		StringData: map[string]string{
			"ARANGO_ROOT_PASSWORD":  "openSesame",
			"DB_CRED":               "root:openSesame",
			"RABBITMQ_USER":         "infinimesh",
			"RABBITMQ_PASS":         "infinimesh",
			"RABBITMQ_CONN":         "amqp://infinimesh:infinimesh@" + instance.Name + "-rabbitmq:5672/",
			"SIGNING_KEY":           "seeeecreet",
			"INF_DEFAULT_ROOT_PASS": "infinimesh",
		},
	}
	if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
		return err
	}
	found := &corev1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, found)
	if errors.IsNotFound(err) {
		return r.Create(context.TODO(), secret)
	}
	return err
}

func (r *ReconcilePlatform) ensurePVC(instance *infinimeshv1beta1.Platform, name string) error {
	qty := resource.MustParse("10Gi")
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: instance.Namespace},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:   corev1.VolumeResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: qty}},
		},
	}
	if err := controllerutil.SetControllerReference(instance, pvc, r.scheme); err != nil {
		return err
	}
	found := &corev1.PersistentVolumeClaim{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: instance.Namespace}, found)
	if errors.IsNotFound(err) {
		return r.Create(context.TODO(), pvc)
	}
	return err
}

func (r *ReconcilePlatform) ensureWorkload(instance *infinimeshv1beta1.Platform, suffix, image string, ports []int, env []corev1.EnvVar, cert bool, dataPath string) error {
	name := instance.Name + "-" + suffix
	var cports []corev1.ContainerPort
	var sports []corev1.ServicePort
	for _, p := range ports {
		proto := corev1.ProtocolTCP
		if suffix == "coap" && p == 5683 {
			proto = corev1.ProtocolUDP
		}
		cports = append(cports, corev1.ContainerPort{ContainerPort: int32(p), Protocol: proto})
		sports = append(sports, corev1.ServicePort{Port: int32(p), Protocol: proto, TargetPort: intstr.FromInt(p)})
	}
	container := corev1.Container{
		Name:            suffix,
		Image:           image,
		ImagePullPolicy: imagePullPolicy(image),
		Env:             env,
		Ports:           cports,
	}
	pod := corev1.PodSpec{Containers: []corev1.Container{container}}
	if suffix == "redis" {
		pod.Containers[0].Args = []string{"redis-server", "--appendonly", "yes"}
	}
	if cert {
		tlsName := instance.Spec.MQTT.SecretName
		if tlsName == "" {
			tlsName = name + "-tls"
		}
		pod.Containers[0].VolumeMounts = append(pod.Containers[0].VolumeMounts, corev1.VolumeMount{Name: "cert", MountPath: "/cert", ReadOnly: true})
		pod.Volumes = append(pod.Volumes, corev1.Volume{Name: "cert", VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: tlsName}}})
	}
	if dataPath != "" {
		pod.Containers[0].VolumeMounts = append(pod.Containers[0].VolumeMounts, corev1.VolumeMount{Name: "data", MountPath: dataPath})
		pod.Volumes = append(pod.Volumes, corev1.Volume{Name: "data", VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: name}}})
	}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: instance.Namespace},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": name}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": name}},
				Spec:       pod,
			},
		},
	}
	if err := r.ensureDeploy(instance, deploy); err != nil {
		return err
	}
	svcType := corev1.ServiceTypeClusterIP
	if suffix == "mqtt-bridge" {
		svcType = corev1.ServiceTypeLoadBalancer
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: instance.Namespace},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Type:     svcType,
			Ports:    sports,
		},
	}
	return r.ensureService(instance, svc)
}

func (r *ReconcilePlatform) ensureDeploy(instance *infinimeshv1beta1.Platform, deploy *appsv1.Deployment) error {
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return err
	}
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if errors.IsNotFound(err) {
		return r.Create(context.TODO(), deploy)
	}
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(deploy.Spec.Template, found.Spec.Template) || !reflect.DeepEqual(deploy.Spec.Selector, found.Spec.Selector) {
		found.Spec.Template = deploy.Spec.Template
		found.Spec.Selector = deploy.Spec.Selector
		return r.Update(context.TODO(), found)
	}
	return nil
}

func (r *ReconcilePlatform) ensureService(instance *infinimeshv1beta1.Platform, svc *corev1.Service) error {
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return err
	}
	found := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found)
	if errors.IsNotFound(err) {
		return r.Create(context.TODO(), svc)
	}
	if err != nil {
		return err
	}
	found.Spec.Selector = svc.Spec.Selector
	found.Spec.Ports = svc.Spec.Ports
	found.Spec.Type = svc.Spec.Type
	return r.Update(context.TODO(), found)
}
