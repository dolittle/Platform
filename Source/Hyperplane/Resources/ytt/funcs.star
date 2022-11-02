metadata = [
        "metadata.annotations[dolittle.io/customer-id]",
        "metadata.annotations[dolittle.io/application-id]",
        "metadata.annotations[dolittle.io/microservice-id]",
        "metadata.labels.customer",
        "metadata.labels.application",
        "metadata.labels.microservice"]

def dolittle_patch_sets_metadata_with_subpath(subPath):
    patches = {x: subPath + x for x in metadata}
    return [{"fromFieldPath": p[0], "toFieldPath": p[1]} for p in patches.items()]
end

def dolittle_patch_sets_metadata():
    patches = dolittle_patch_sets_metadata_with_subpath("") + dolittle_patch_sets_metadata_with_subpath("spec.forProvider.manifest.")
    return [{"name": "dolittle-metadata", "patches": patches}]
end

def dolittle_references(apiVersion, kind, fromFieldPath, namespace="",toFieldPaths=[]):
    toFieldPaths = ["", "spec.forProvider.manifest."] + toFieldPaths
    patches = [{"fromFieldPath": fromFieldPath, "toFieldPath": x} for x in ["spec.references[" + str(i) +"].patchesFrom.name" for i in range(len(toFieldPaths))]]
    references = []
    for toFieldPath in toFieldPaths:
        for fieldPath in metadata:
            item = {"toFieldPath": toFieldPath + fieldPath}
            reference = {"apiVersion": apiVersion, "kind": kind, "fieldPath": fieldPath}
            if namespace != "":
                reference["namespace"] = namespace
            end
            item["patchesFrom"] = reference
            references.append(item)
        end
    end
    return {
        "patchSets": [{
            "name": "dolittle-references",
            "patches": patches
        }],
        "references": references
    }
end
