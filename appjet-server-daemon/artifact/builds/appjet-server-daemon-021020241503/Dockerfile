FROM --platform=linux/arm64 maven:3.8.2-openjdk-17-slim as builder

# Install builder dependencies
RUN apt-get update && apt-get install -y git

WORKDIR /app

# Code checkout
ARG GIT_USERNAME=guilhermemalhado1
ARG GIT_PASSWORD=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
RUN git config --global credential.helper '!f() { echo "username=${GIT_USERNAME}"; echo "password=${GIT_PASSWORD}"; }; f'
RUN git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/guilhermemalhado1/javadummyproj.git /app || (echo "Git clone failed"; exit 1)


COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Build the application
RUN mvn install -DskipTests

FROM --platform=linux/arm64 arm64v8/openjdk:17-jdk-slim-buster

WORKDIR /app

RUN apt-get update && apt-get install -y netcat

COPY --from=builder /wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

COPY --from=builder /app/target/demo-0.0.1-SNAPSHOT.jar app.jar

EXPOSE 8080

CMD ["/wait-for-it.sh", "database:3306", "--", "java", "-jar", "app.jar"]
