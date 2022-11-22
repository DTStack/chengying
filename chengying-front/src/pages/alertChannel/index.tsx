import { asyncComponent } from 'react-async-component';

export default asyncComponent({
  name: 'AlertChannel',
  serverMode: 'resolve',
  resolve: async () => {
    const module = await import('./alertChannel');
    return module.default;
  },
});
