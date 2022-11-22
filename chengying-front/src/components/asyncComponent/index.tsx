import { asyncComponent } from 'react-async-component';
import LoadingComponent from '@/components/loading/loading';
import ErrorComponent from '@/components/error/error';

export default (component, name) =>
  asyncComponent({
    name: name,
    serverMode: 'resolve',
    resolve: async () => {
      const module = await component;
      return module.default;
    },
    LoadingComponent: LoadingComponent,
    ErrorComponent: ErrorComponent,
  });
