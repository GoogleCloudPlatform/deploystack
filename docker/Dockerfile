FROM golang:latest as builder

WORKDIR /app/ds

COPY deploystack_exec.go deploystack_exec.go 
COPY version.sh version.sh 
RUN go mod init exec 
RUN go get github.com/GoogleCloudPlatform/deploystack/config
RUN go get github.com/GoogleCloudPlatform/deploystack/tui
RUN go get github.com/GoogleCloudPlatform/deploystack/gcloud
RUN go get github.com/GoogleCloudPlatform/deploystack
RUN go mod tidy
RUN ./version.sh
RUN GOOS=linux GOARCH=amd64 go build -o dsexec deploystack_exec.go

FROM google/cloud-sdk:latest as gcloud
COPY deploystack /usr/local/bin/deploystack
COPY --from=builder /app/ds/dsexec /usr/local/bin/dsexec

RUN apt-get update && apt-get install -y \
    wget unzip curl \
  && rm -rf /var/lib/apt/lists/*

ARG TERRAFORM_VERSION=1.3.9
RUN wget --quiet https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && mv terraform /usr/bin \
  && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip
