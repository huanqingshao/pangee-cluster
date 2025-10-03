---
---

# 重置密码

本文描述了如何重置 PangeeCluster 中 `admin` 用户的密码。请按照如下步骤操作：

- 执行如下命令，进入 pangee-cluster 容器的命令行界面：

  ```sh
  docker exec -it pangee-cluster /bin/bash
  ```

- 在 pangee-cluster 容器中执行如下命令：

  ```sh
  ./pangee-cluster-admin reset-password
  ```

  输出结果如下所示：

  ```log
  设置日志级别为 info
  2022/04/30 22:41:39.801 |  info | try to reset-password
  2022/04/30 22:41:39.803 |  info | 已将 admin 的密码重置为 PangeeCluster123
  ```

  至此，您已经成功重置了 `admin` 用户的密码，重置后密码为 `PangeeCluster123`
