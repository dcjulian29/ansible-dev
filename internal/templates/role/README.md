# Ansible Role: !!ROLE_NAME!!

[![Lint](https://github.com/dcjulian29/ansible-role-!!ROLE_NAME!!/actions/workflows/lint.yml/badge.svg)](https://github.com/dcjulian29/ansible-role-!!ROLE_NAME!!/actions/workflows/lint.yml) [![GitHub Issues](https://img.shields.io/github/issues-raw/dcjulian29/ansible-role-!!ROLE_NAME!!.svg)](https://github.com/dcjulian29/ansible-role-!!ROLE_NAME!!/issues)

This an Ansible role to !!ROLE_DESC!!

## Requirements

- Active Internet Connection.

## Installation

To use, use `requirements.yml` with the following git source:

```yaml
---
roles:
- name: dcjulian29.!!ROLE_NAME!!
  src: https://github.com/dcjulian29/ansible-role-!!ROLE_NAME!!.git
  version: main
  ```

Then download it with `ansible-galaxy`:

```shell
ansible-galaxy install -r requirements.yml
```

## Dependencies

- None
