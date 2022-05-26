FROM golang:1.17-alpine3.14 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED 0
RUN go build ./cmd/give-me-bnb

FROM python:3.9-slim
RUN apt update && apt install -y tor netcat wget && \
    wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt install -y ./google-chrome-stable_current_amd64.deb && \
    rm google-chrome-stable_current_amd64.deb && \
    apt clean
WORKDIR /app
COPY ./third_party/hcaptcha-challenger/requirements.txt /app/third_party/hcaptcha-challenger/requirements.txt
RUN pip3 install -r /app/third_party/hcaptcha-challenger/requirements.txt
COPY ./third_party /app/third_party
RUN python3 /app/third_party/hcaptcha-challenger/src/main.py install
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
COPY --from=builder /app/give-me-bnb .
ENV PATH="${PWD}:${PATH}"
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["give-me-bnb"]
