if ('serviceWorker' in navigator) {
  window.addEventListener('load', function () {
    navigator.serviceWorker
      .register('sw.js', { scope: '/' })
      .then(function (registration) {
        // 注册成功
        console.log(
          'ServiceWorker registration successful with scope: ',
          registration.scope
        );
      })
      .catch(function (err) {
        // 注册失败:(
        console.log('ServiceWorker registration failed: ', err);
      });
  });
}
