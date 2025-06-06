export default function () {
  if (window.PangeeCluster.version.arch === 'arm64') {
    return "https://addons.kuboard.cn/v-pangee-cluster-arm64"
  }
  return "https://addons.kuboard.cn/v-pangee-cluster"
}