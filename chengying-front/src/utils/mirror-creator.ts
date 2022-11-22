export default (items: any, options: any = {}) => {
  if (!Array.isArray(items)) {
    throw new Error('mirrorCreator(...): argument must be an array.');
  }

  const { prefix } = options;
  const container = {};
  items.forEach((item) => (container[item] = `${prefix || ''}${item}`));
  return container;
};
