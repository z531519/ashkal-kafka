FROM gradle:jdk17 as base
COPY . .
RUN gradle assemble

RUN pwd; ls -l build/libs

FROM gradle:jdk17
# copy jar from build image
COPY --from=base /home/gradle/build/libs/kafka-templates-0.0.1.jar /app/demoapp.jar

ENTRYPOINT ["java"]
CMD ["-jar", "/app/demoapp.jar"]


