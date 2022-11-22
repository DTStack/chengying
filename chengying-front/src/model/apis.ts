export default interface ApiModel {
  sort: string;
}

export interface ServiceTree {
  productName: string;
  productVersion: number;
  serviceName: string;
}

export interface ServiceFile {
  productName: string;
  productVersion: number;
  serviceName: string;
  file: string;
}

export interface DistributeServiceConfig {
  serviceName: string;
  productId: number;
}

export interface SearchDeployLogs {
  deployId: string;
  serviceName: string;
  productName: string;
  productVersion: number;
}
