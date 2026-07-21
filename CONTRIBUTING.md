# Contributing to MISP-Operator

## Development

The MISP-Operator was developed using the Kubebuilder project.

The easiest way to start developing in this project is to [spin up a local `kind` cluster](https://book.kubebuilder.io/reference/kind.html).
Make sure it is selected using `KUBECONFIG`.

#### Commit Message Convention

The project uses [Release Please](https://github.com/googleapis/release-please) for GitHub release management and [CHANGELOG](CHANGELOG.md) generation. Therefore it is necessary to follow the [Conventional Commits](https://www.conventionalcommits.org) convention.

In summary this means a commit message format of:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Where the `type` is one of `fix:`, `feat:`, `refactor:`, `chore:`, `docs:`, `perf:`.

Breaking changes need to be highlighted using a `!` after the type/scope or by putting `BREAKING CHANGE` somewhere in the body of the commit message:

```
fix!: bla bla bla
```

```
fix: bla bla bla

BREAKING CHANGE: This breaks x because of yz.
```

#### Running locally

To develop things locally you can use the `make run` task for a simple and fast way to run the operator in the foreground. This is great for quick feedback and code-level debugging.

#### Package the operator

You can also package the operator and load it into a local development cluster:

```shell
make docker-build IMG=ghcr.io/pascaliske/misp-operator:main
make deploy IMG=ghcr.io/pascaliske/misp-operator:main
```

#### Generating CRD Manifests

To generate and install the manifests of all CRDs in a local cluster you can use the following tasks:

```shell
make manifests # generate the manifests in config/crd/bases
make install # install the manifests in the currently configured cluster
```

> [!TIP]
> Run `make help` for more information on all potential `make` targets

More information can be found in the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html).
