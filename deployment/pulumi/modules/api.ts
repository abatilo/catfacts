import * as pulumi from "@pulumi/pulumi";
import * as awsx from "@pulumi/awsx";
import * as k8s from "@pulumi/kubernetes";

interface Args {
  k8sNamespace: pulumi.Input<string>;
}

class API extends pulumi.ComponentResource {
  public constructor(
    name: string,
    { k8sNamespace }: Args,
    opts?: pulumi.ComponentResourceOptions
  ) {
    super("abatilo:catfacts:api", name, {}, opts);
    const defaultOptions: pulumi.CustomResourceOptions = { parent: this };

    const config = new pulumi.Config();

    const postgres = new k8s.helm.v3.Chart(
      `postgres`,
      {
        namespace: k8sNamespace,
        fetchOpts: {
          repo: "https://charts.bitnami.com/bitnami",
        },
        chart: "postgresql",
        version: "10.3.11",
        values: {
          global: { storageClass: "gp2" },
          image: { tag: "9.6.12" },
          postgresqlPassword: config.requireSecret("postgresPassword"),
          postgresqlPostgresPassword: config.requireSecret("postgresPassword"),
          rbac: { create: true },
        },
      },
      defaultOptions
    );

    const repository = new awsx.ecr.Repository(name, {}, defaultOptions);
    const image = repository.buildAndPushImage({
      context: "../../",
      dockerfile: "../../build/Dockerfile.api",
      cacheFrom: {
        stages: ["build"],
      },
      env: { DOCKER_BUILDKIT: "1" },
    });

    const api = new k8s.helm.v3.Chart(
      name,
      {
        path: "../helm/api",
        namespace: k8sNamespace,
        values: {
          replicaCount: 2,
          initContainers: {},
          image: { repository: image },
          ingress: {
            match: "Host(`catfacts.aaronbatilo.dev`) && PathPrefix(`/api`)",
          },
          env: {
            CF_TWILIO_ACCOUNT_SID: config.requireSecret("twilio_account_sid"),
            CF_TWILIO_AUTH_TOKEN: config.requireSecret("twilio_auth_token"),
            CF_TWILIO_PHONE_NUMBER: config.requireSecret("twilio_phone_number"),
            CF_DB_HOST: "postgres-postgresql",
            CF_DB_PASSWORD: config.requireSecret("postgresPassword"),
          },
          entryPoints: ["websecure"],
        },
      },
      { ...defaultOptions, dependsOn: [postgres] }
    );

    // Create an example CronJob.
    const cronBlast = new k8s.batch.v1beta1.CronJob(
      name,
      {
        spec: {
          schedule: "10 18 * * *",
          jobTemplate: {
            spec: {
              template: {
                spec: {
                  containers: [
                    {
                      env: [
                        {
                          name: "CF_TWILIO_ACCOUNT_SID",
                          value: config.requireSecret("twilio_account_sid"),
                        },
                        {
                          name: "CF_TWILIO_AUTH_TOKEN",
                          value: config.requireSecret("twilio_auth_token"),
                        },
                        {
                          name: "CF_TWILIO_PHONE_NUMBER",
                          value: config.requireSecret("twilio_phone_number"),
                        },
                        {
                          name: "CF_DB_HOST",
                          value: "postgres-postgresql",
                        },
                        {
                          name: "CF_DB_PASSWORD",
                          value: config.requireSecret("postgresPassword"),
                        },
                      ],
                      name: "blast",
                      image,
                      command: ["cf"],
                      args: ["blast"],
                    },
                  ],
                  restartPolicy: "Never",
                },
              },
            },
          },
        },
      },
      { ...defaultOptions, dependsOn: [postgres, api] }
    );
  }
}
export default API;
