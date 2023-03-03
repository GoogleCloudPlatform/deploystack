FROM golang:latest as builder

ARG DS_BASE_VERSION=@latest
ARG DS_GITHUB_VERSION=@latest
ARG DS_TUI_VERSION=@latest
ARG DS_GCLOUD_VERSION=@latest
ARG DS_CONFIG_VERSION=@latest

WORKDIR /app/ds

COPY deploystack_exec.go deploystack_exec.go 
COPY version.sh version.sh 
RUN go mod init exec
RUN go get github.com/GoogleCloudPlatform/deploystack/github${DS_GITHUB_VERSION}
RUN go get github.com/GoogleCloudPlatform/deploystack/config${DS_CONFIG_VERSION}
RUN go get github.com/GoogleCloudPlatform/deploystack/tui${DS_TUI_VERSION}
RUN go get github.com/GoogleCloudPlatform/deploystack/gcloud${DS_GCLOUD_VERSION}
RUN go get github.com/GoogleCloudPlatform/deploystack${DS_BASE_VERSION}
RUN go mod tidy
RUN ./version.sh
RUN GOOS=linux GOARCH=amd64 go build -o dsexec deploystack_exec.go

FROM google/cloud-sdk:latest as gcloud

ARG TERRAFORM_VERSION=1.3.9

COPY deploystack /usr/local/bin/deploystack
COPY --from=builder /app/ds/dsexec /usr/local/bin/dsexec

RUN apt-get update && apt-get install -y \
    wget unzip curl \
  && rm -rf /var/lib/apt/lists/*


RUN wget --quiet https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
  && mv terraform /usr/bin \
  && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip
