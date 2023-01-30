FROM alpine:latest
COPY worker/ /app/worker/
WORKDIR /app/worker
EXPOSE 7777
CMD ["./worker"]