# Ansible Development Tool (ansible-dev)

[![Version](https://img.shields.io/github/v/release/dcjulian29/ansible-dev)](https://github.com/dcjulian29/ansible-dev/releases)
[![GitHub Issues](https://img.shields.io/github/issues-raw/dcjulian29/ansible-dev.svg)](https://github.com/dcjulian29/ansible-dev/issues)
[![Build](https://github.com/dcjulian29/ansible-dev/actions/workflows/build.yml/badge.svg)](https://github.com/dcjulian29/ansible-dev/actions/workflows/build.yml)


ansible-dev integrates with Vagrant to enable developers to define, develop, and test Ansible
playbooks, roles, and runbooks.

By utilizing Ansible, developers/operators can automate the deployment of software applications
across multiple hosting providers, reducing the time and effort required to manage complex
infrastructure environments.

## Configuration

ansible-dev reads its settings from `~/.config/ansible-dev.yml`. Run `ansible-dev config path` to
print the exact location or `ansible-dev config show` to dump the current values. A fully
commented file you can copy is included as
[ansible-dev.yml.example](ansible-dev.yml.example).

Two settings are required before the role and runbook commands will work:

| Setting | Purpose |
| --- | --- |
| `roles_path` | Directory holding your published role repositories. Required by `role compare` and `role new --publish`. |
| `runbooks_path` | Directory holding your published runbook repositories. Required by `runbook compare` and `runbook new`. |

```shell
ansible-dev config roles-path /path/to/ansible/roles
ansible-dev config runbooks-path /path/to/ansible/runbooks
```

Both expect an absolute directory. A relative path, or one that does not exist, is still saved but
reported as a warning — so a configuration file can be prepared on a machine where the directory
has yet to be created, without a typo passing unnoticed.

### Upgrading from the environment variables

Earlier versions read the paths from `ANSIBLE_ROLES` and `ANSIBLE_RUNBOOKS`. Those variables are
no longer consulted. To carry existing values into the configuration file:

```shell
ansible-dev config import-env
```

Settings that already have a value are left unchanged; pass `--force` to overwrite them. Any
command that needs an unset path will point you at `config import-env` when the matching variable
is still present in your environment.

### Compare ignore lists

`role compare` and `runbook compare` skip any path containing one of the configured substrings.
An empty list compares everything.

Each entry is a plain substring test against the whole path — not a regular expression and not a
glob — so a short entry excludes more than it might appear to. `.git` also excludes `.github`,
`.gitignore` and `.gitattributes`; `.ansible` also excludes `.ansible-lint`. For the same reason
entries should not contain escapes or path separators: `\.git` matches only on Windows, where the
backslash happens to be the separator, and silently matches nothing on Linux or macOS.

```shell
ansible-dev config role-ignore list
ansible-dev config role-ignore add .vscode
ansible-dev config role-ignore remove .vscode
ansible-dev config role-ignore clear
```

`runbook-ignore` accepts the same subcommands for the runbook list.

### Diff tool

The external diff tool is configured per operating system, so one file can serve several machines
and supporting a new platform is just a new entry. Only the entry for the host you are running on
is consulted, and it is looked up only when a diff is about to open — so a missing diff program is
reported only when something actually differs, never on a clean comparison and never under
`--no-diff`.

```shell
ansible-dev config diff-program "C:\Program Files\WinMerge\winmergeu.exe"
ansible-dev config diff-role-filter AnsibleRoles
ansible-dev config diff-runbook-filter AnsibleRunbooks
ansible-dev config diff-args /r /m Full /u /f "{filter}" "{left}" "{right}"
```

The arguments are a template: `{left}` and `{right}` are replaced with the two directories being
compared, and `{filter}` with the role or runbook filter for the comparison being run. When that
filter is empty, a standalone `{filter}` argument is dropped along with the flag immediately
before it, so no dangling flag is passed to the diff tool.

The WinMerge example above uses `/`-style switches, which pass through untouched. A tool that
takes Unix-style flags needs `--` first, otherwise `ansible-dev` treats a leading `-` as one of
its own flags and fails with `unknown shorthand flag`:

```shell
ansible-dev config diff-args -- -r --brief "{left}" "{right}"
```
