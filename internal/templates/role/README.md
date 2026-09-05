# Ansible Role: !!ROLE_NAME!!

[![Lint](https://github.com/!!ROLE_NAMESPACE!!/ansible-role-!!ROLE_NAME!!/actions/workflows/lint.yml/badge.svg)](https://github.com/!!ROLE_NAMESPACE!!/ansible-role-!!ROLE_NAME!!/actions/workflows/lint.yml) [![GitHub Issues](https://img.shields.io/github/issues-raw/!!ROLE_NAMESPACE!!/ansible-role-!!ROLE_NAME!!.svg)](https://github.com/!!ROLE_NAMESPACE!!/ansible-role-!!ROLE_NAME!!/issues)

!!ROLE_DESC!!

## Requirements

- Active Internet Connection.

## Installation

To use, use `requirements.yml` with the following git source:

```yaml
---
roles:
- name: !!ROLE_NAMESPACE!!.!!ROLE_NAME!!
  src: https://github.com/!!ROLE_NAMESPACE!!/ansible-role-!!ROLE_NAME!!.git
  version: main
  ```

Then download it with `ansible-galaxy`:

```shell
ansible-galaxy install -r requirements.yml
```

## Dependencies

- None
