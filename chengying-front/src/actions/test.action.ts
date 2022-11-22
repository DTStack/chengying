export const changeUsername = (name: string) => {
  return {
    type: 'CHANGE_USERNAME',
    payload: name,
  };
};

export interface TestActionTypes {
  changeUsername: Function;
}
