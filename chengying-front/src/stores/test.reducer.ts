interface ActionType {
  type: string;
  payload?: any;
}

interface TestStore {
  user: {
    username: string;
    sex: string;
    age: number;
  };
}

export type TestStoreTypes = TestStore;

const initState: TestStore = {
  user: {
    username: '',
    sex: '',
    age: 20,
  },
};

export default (state = initState, action: ActionType) => {
  const { type, payload } = action;
  switch (type) {
    case 'CHANGE_USERNAME':
      const user = { ...state.user, username: payload };
      return { ...state, user: user };
    default:
      return state;
  }
};
