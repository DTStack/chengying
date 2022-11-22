import { asyncComponent } from 'react-async-component';

export default asyncComponent({
  name: 'ServicePage',
  serverMode: 'resolve',
  resolve: async () => {
    const module = await import('./container');
    return module.default;
  },
});
