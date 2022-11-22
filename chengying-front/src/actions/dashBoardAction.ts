import { Dispatch } from 'redux';
import { dashboardService } from '@/services';
import { message } from 'antd';
import { DashBoardActions } from '@/constants/actionTypes';
import * as _ from 'lodash';

// interface dashEntity {
//   title: String,
//   id: number,
//   list: Array<any>
// }
export const getDashboardList = (params: object) => {
  return (dispatch: Dispatch) => {
    dashboardService
      .getServiceDashInfo(params)
      .then((res: any) => {
        let tags: any = [];
        const dash: any = [
          {
            title: 'General',
            id: 0,
            list: [],
          },
        ];
        for (const d of res.data) {
          if (d.tags.length) {
            tags = _.concat(tags, d.tags);
          }
          switch (d.type) {
            case 'dash-folder':
              d.list = [];
              dash.push(d);
              break;
            case 'dash-db':
              if (d.folderId) {
                for (const s in dash) {
                  if (d.folderId === dash[s].id) {
                    dash[s].list.push(d);
                  }
                }
              } else {
                dash[0].list.push(d);
              }
              break;
          }
        }
        dispatch({
          type: DashBoardActions.UPDATE_DASH_LIST,
          payload: {
            dashboards: dash,
            tags: _.uniq(tags),
          },
        });
      })
      .catch((err: any) => message.error(err));
  };
};

export interface DashBoardActionsTypes {
  getDashboardList: Function;
}
