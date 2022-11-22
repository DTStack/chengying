import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: '自动化部署',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
      <>
        ChengYing通过规范化的步骤和参数约定制作出产品安装包，发布包中的Schema文件中配置了安装包中所有的服务，包含各服务的配置参数、健康检查参数，服务之间的依赖关系等。产品部署时可根据Schema中的相关配置实现一键全自动化部署。
      </>
    ),
  },
  {
    title: '界面化集群运维',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
      <>
        Hadoop集群、大数据平台在日常运维中涉及到的节点扩容缩容、组件停止启动、服务滚动重启、服务参数修改、版本升级与回滚等多种运维操作，通过逻辑化、流程化的产品界面展现，方便运维人员操作和监控，提高运维效率。
      </>
    ),
  },
  {
    title: '仪表盘集群监控',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
      <>
        通过集成开源的promethus和grafana，实现对集群、服务、节点的核心参数监控，并通过灵活形象的仪表盘进行数据展现。包含CPU占用率，RAM使用率、磁盘空间、I/O读写速率等核心参数进行监控，实时掌握集群、服务、节点的运行状态，降低运维故障率。同时，支持用户自建仪表盘及监控项，实现自定义监控项。
      </>
    ),
  },
  {
    title: '实时告警',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
        <>
          支持实时监控集群中各组件服务的运行指标，如CPU、内存、磁盘、读写IO等，并支持短信、钉钉、邮件告警通道配置，集成多种第三方消息插件。当集群服务出现异常时，可触发告警条件，系统将及时通知接收人。
        </>
    ),
  },
  {
    title: '强扩展性',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
        <>
          通过自研的Agent Server抽象出七大REST接口，安装、启动、停止、更新、配置修改、卸载、执行等与上层应用进行交互，可使agent类别和功能可轻松无限扩展。
        </>
    ),
  },
  {
    title: '安全稳定',
    Svg: require('@site/static/img/chengying-dark.svg').default,
    description: (
        <>
          数据安全、产品安全是大数据产品需要重点考虑的问题。ChengYing在产品设计中过滤掉rm、drop等命令行，防止对数据库的误操作，通过更加安全的方式执行相关命令。同时提供服务的滚动重启、产品的断电重启，解决运维时服务不停止运行的场景并节省运维时间。
        </>
    ),
  }
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
