# Ansible Runbook: !!RUNBOOK_NAME!!

[![GitHub Issues](https://img.shields.io/github/issues-raw/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!.svg)](https://github.com/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!/issues)
[![Version](https://img.shields.io/github/v/release/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!)](https://github.com/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!/releases)
[![Build](https://github.com/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!/actions/workflows/build.yml/badge.svg)](https://github.com/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!/actions/workflows/build.yml)

This is an Ansible runbook that will !!RUNBOOK_DESC!!

## Requirements

- Active Internet Connection.

## Installation

To use, use `requirements.yml` with the following git source:

```yaml
---
collections:
- name: dcjulian29.!!RUNBOOK_NAME!!
  type: git
  source: https://github.com/dcjulian29/ansible-runbook-!!RUNBOOK_NAME!!.git
  ```

Then download it with `ansible-galaxy`:

```shell
ansible-galaxy collection install -r requirements.yml
```

To excute the runbook:

```shell
ansible-playbook dcjulian29.!!RUNBOOK_NAME!!.runbook.yml
```
