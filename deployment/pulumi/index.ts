import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import API from "./modules/api";
import Frontend from "./modules/frontend";

if (pulumi.getStack() === "prod") {
  const cfNamespace = new k8s.core.v1.Namespace("catfacts", {
    metadata: {
      name: "catfacts",
    },
  });

  const api = new API("api", { k8sNamespace: cfNamespace.metadata.name });
  const frontend = new Frontend("frontend", {
    k8sNamespace: cfNamespace.metadata.name,
  });
}
