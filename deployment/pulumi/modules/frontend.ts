import * as pulumi from "@pulumi/pulumi";
import * as awsx from "@pulumi/awsx";
import * as k8s from "@pulumi/kubernetes";

interface Args {
  k8sNamespace: pulumi.Input<string>;
}

class Frontend extends pulumi.ComponentResource {
  public constructor(
    name: string,
    { k8sNamespace }: Args,
    opts?: pulumi.ComponentResourceOptions
  ) {
    super("abatilo:catfacts:frontend", name, {}, opts);
    const defaultOptions: pulumi.CustomResourceOptions = { parent: this };

    const repository = new awsx.ecr.Repository(name, {}, defaultOptions);
    const image = repository.buildAndPushImage({
      context: "../../web/frontend",
      dockerfile: "../../build/Dockerfile.frontend",
      cacheFrom: {
        stages: ["build"],
      },
      env: { DOCKER_BUILDKIT: "1" },
    });

    const api = new k8s.helm.v3.Chart(
      name,
      {
        path: "../helm/frontend",
        namespace: k8sNamespace,
        values: {
          replicaCount: 2,
          image: { repository: image },
          ingress: {
            match: "Host(`catfacts.aaronbatilo.dev`) && PathPrefix(`/`)",
          },
          entryPoints: ["websecure"],
        },
      },
      defaultOptions
    );
  }
}
export default Frontend;
