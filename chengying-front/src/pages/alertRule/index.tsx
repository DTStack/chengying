import { asyncComponent } from 'react-async-component';
import LoadingComponent from '@/components/loading/loading';
import ErrorComponent from '@/components/error/error';

export default asyncComponent({
  name: 'AlertRulePage',
  resolve: async () => {
    const module = await import('./rule');
    return module.default;
  },
  LoadingComponent: LoadingComponent,
  ErrorComponent: ErrorComponent,
});
