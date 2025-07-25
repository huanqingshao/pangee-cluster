export default function (source) {
  if (source === "github") {
    return "https://raw.githubusercontent.com/Horan-Z/test/refs/heads"
  }
  if (source === "gitee") {
    return "https://gitee.com/Horan-Z/test/raw"
  }
}