# Avoiding security risks with ansible.cfg in the current directory

[WARNING]: Ansible is being run in a world writable directory (/workspaces/pangee-cluster-resource-package), ignoring it as an ansible.cfg source. For more information see
https://docs.ansible.com/ansible/devel/reference_appendices/config.html#cfg-in-world-writable-dir
ERROR! the role 'init-os' was not found in /workspaces/pangee-cluster-resource-package/operations/install-cluster/01-init-os/roles:/root/.ansible/roles:/usr/share/ansible/roles:/etc/ansible/roles:/workspaces/pangee-cluster-resource-package/operations/install-cluster/01-init-os

The error appears to be in '/workspaces/pangee-cluster-resource-package/operations/install-cluster/01-init-os/playbook.yaml': line 6, column 7, but may
be elsewhere in the file depending on the exact syntax problem.

The offending line appears to be:

  roles:
    - init-os
      ^ here