export default interface DeployLog {
  id: number;
  productName: string;
  productVersion: number | string;
  service?: string;
}
