export default [
    {
      path: '/maxkb',
      children: [
        {
          path: '',
          component: () => import('@/views/maxkb/index.vue'),
          meta: {
            title: '客服应用管理',
            activePath: '/maxkb'
          }
        },
    ]
}
]