import apis from '@/constants/apis';
import * as http from '@/utils/http';

const { backup } = apis

export default {
  queryBuildBackupPath(params: any) {
    return http[backup.queryBuildBackupPath.method](
      backup.queryBuildBackupPath.url,
      params
    );
  },
  SetUpBackupPath(params: any) {
    return http[backup.SetUpBackupPath.method](
      backup.SetUpBackupPath.url,
      params
    )
  }
};