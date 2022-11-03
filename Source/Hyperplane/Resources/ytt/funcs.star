metadata = [
        "metadata.annotations[dolittle.io/customer-id]",
        "metadata.annotations[dolittle.io/tenant-id]",
        "metadata.annotations[dolittle.io/application-id]",
        "metadata.annotations[dolittle.io/microservice-id]",
        "metadata.labels.customer",
        "metadata.labels.application",
        "metadata.labels.microservice"]
labels = [
        "customer-id",
        "application-id",
        "microservice-id",
        "customer",
        "application",
        "microservice",
        "environment"]

metadataFieldSubPaths = {
    "customer-id": "metadata.annotations[dolittle.io/customer-id]",
    "metadata.annotations[dolittle.io/customer-id]": "metadata.annotations[dolittle.io/customer-id]",
    "metadata.annotations[dolittle.io/tenant-id]": "metadata.annotations[dolittle.io/customer-id]",
    "application-id": "metadata.annotations[dolittle.io/application-id]",
    "metadata.annotations[dolittle.io/application-id]": "metadata.annotations[dolittle.io/application-id]",
    "microservice-id": "metadata.annotations[dolittle.io/microservice-id]",
    "metadata.annotations[dolittle.io/microservice-id]": "metadata.annotations[dolittle.io/microservice-id]",
    "customer": "metadata.labels.customer",
    "metadata.labels.customer": "metadata.labels.customer",
    "tenant": "metadata.labels.customer",
    "metadata.labels.tenant": "metadata.labels.customer",
    "application": "metadata.labels.application",
    "metadata.labels.application": "metadata.labels.application",
    "microservice": "metadata.labels.microservice",
    "metadata.labels.microservice": "metadata.labels.microservice",
    "environment": "metadata.labels.environment",
    "metadata.labels.environment": "metadata.labels.environment"}

def dolittle_patch_sets_metadata_with_subpath(dolittleMetadata, fromFieldSubPath, toFieldSubPath):
    patches = {fromFieldSubPath + x: toFieldSubPath + metadataFieldSubPaths[x] for x in dolittleMetadata}
    return [{"fromFieldPath": p[0], "toFieldPath": p[1]} for p in patches.items()]
end

def dolittle_patch_sets_metadata_from_metadata():
    patches = dolittle_patch_sets_metadata_with_subpath(metadata, "", "")
    patches = patches + dolittle_patch_sets_metadata_with_subpath(metadata, "", "spec.forProvider.manifest.")
    return [{"name": "dolittle-metadata", "patches": patches}]
end

def dolittle_patch_sets_metadata_from_spec(specPath="spec.metadata."):
    patches = dolittle_patch_sets_metadata_with_subpath(labels, specPath, "")
    patches = patches + dolittle_patch_sets_metadata_with_subpath(labels, specPath, "spec.forProvider.manifest.")
    return [{"name": "dolittle-metadata", "patches": patches}]
end

def dolittle_references(apiVersion, kind, fromFieldPath, namespace="",toFieldPaths=[]):
    toFieldPaths = ["", "spec.forProvider.manifest."] + toFieldPaths
    references = []
    for toFieldPath in toFieldPaths:
        for fieldPath in metadata:
            item = {"toFieldPath": toFieldPath + metadataFieldSubPaths[fieldPath]}
            reference = {"apiVersion": apiVersion, "kind": kind, "fieldPath": fieldPath}
            if namespace != "":
                reference["namespace"] = namespace
            end
            item["patchesFrom"] = reference
            references.append(item)
        end
    end

    patches = [{"fromFieldPath": fromFieldPath, "toFieldPath": x} for x in ["spec.references[" + str(i) +"].patchesFrom.name" for i in range(len(references))]]
    return {
        "patchSets": [{
            "name": "dolittle-references",
            "patches": patches
        }],
        "references": references
    }
end
