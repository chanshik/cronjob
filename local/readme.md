# 로컬 환경에서 실행
## 필요한 작업

### 1. 기존에 배포되어 있는 Webhook 설정 & 컨트롤러 컨테이너 제거

```shell
make undeploy
```

### 2. 보안 채널을 생성할 때 사용할 인증서와 비밀키 생성

Webhook 요청을 받을 머신의 IP 주소 혹은 DNS 이름을 넣어서 생성해야 접속할 때 사용할 수 있다.

```shell
openssl req -new -x509 \
  -subj "/CN=webhook.cronjob.svc" -nodes \
  -newkey rsa:4096 -keyout tls.key \
  -out tls.crt -days 365 \
  -addext "subjectAltName = DNS:webhook.cronjob.svc,DNS:webhook.cronjob.svc.cluster.local,IP:192.168.50.193"
```

### 3. 생성한 인증서를 `base64` 인코딩을 통해 변환

```shell
base64 tls.crt -w 0
```

### 4. Webhook 설정 안에 `clientConfig` 설정 변경

기존 webhook 설정에는 webhook 요청을 보낼 `service` 객체를 지정하였지만,
로컬 머신에서 실행하기 위해 IP 주소와 API 경로를 입력한다.

Kubernetes 클러스터에서 로컬 머신을 직접 접근할 수 없을 경우에는 
`ngrok` 서비스를 이용해 로컬 머신으로 접근할 수 있는 경로를 열 수 있다.

`caBundle` 항목에는 인증서를 `base64` 인코딩으로 변환한 값을 넣는다. 

```yaml
...
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  creationTimestamp: null
  name: mutating-webhook-configuration
webhooks:
  - admissionReviewVersions:
      - v1
      - v1beta1
    clientConfig:
      url: https://192.168.50.193:9443/mutate-batch-tutorial-chanshik-dev-v1-cronjob
      caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FUR.....DQVRFLS0tLS0K
...

---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  creationTimestamp: null
  name: validating-webhook-configuration
webhooks:
  - admissionReviewVersions:
      - v1
      - v1beta1
    clientConfig:
      url: https://192.168.50.193:9443/validate-batch-tutorial-chanshik-dev-v1-cronjob
      caBundle: LS0tLS1CRUdJTiBDR.......klDQVRFLS0tLS0K
...
```

### 5. CronJob CRD 와 로컬 머신 주소를 기록한 Webhook 설정 배포

Tutorial 1

```shell
kubectl apply -f cronjob-crd.yaml
kubectl apply -f webhook-manifests.yaml
```

Tutorial 2

```shell
kubectl apply -f cronjob-crd-multiversion.yaml
kubectl apply -f webhook-manifests.yaml
```
