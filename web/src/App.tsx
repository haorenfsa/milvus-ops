import React from 'react';
import Home from './pages/Home'
import './App.css';

import { Router, Route, Switch, Redirect } from 'react-router';
import { createBrowserHistory } from 'history'
import { Layout, Menu, Icon } from 'antd';
import { Link } from 'react-router-dom';
import Shell from './pages/Shell';
import Log from './pages/Log';


const { Header, Content, Footer, Sider } = Layout;

interface MenuData{
  path: string
  icon: string
  name: string
  component: any
}

const menudatas: MenuData[] = [
  {
    path: "/app/",
    name: "Dashboard",
    icon: "dashboard",
    component: Home,
  },
  {
    path: "/app/shell",
    name: "WebShell",
    icon: "code",
    component: Shell,
  },
  {
    path: "/app/logs",
    name: "Logs",
    icon: "code",
    component: Log,
  }
]

class App extends React.Component {
  state = {
    collapsed: false,
  }
  render() {
    const { collapsed } = this.state
    const history = createBrowserHistory()
    const path = history.location.pathname
    let currentMenuKey = '1'
    menudatas.forEach((menudata: MenuData, index: number) => {
      if (menudata.path === path) {
        currentMenuKey = `${index+1}`
      }
    })
    return (
      <div className="App">
        <Router history={history} >
          <Layout style={{ minHeight: '100vh' }}>
            <Sider collapsible collapsed={collapsed} onCollapse={(collapsed: boolean) => this.setState({collapsed})}>
              {/* <div className='logo'>
                <img style={{width: "100%"}} src={'/logo.png'} alt="" />
              </div> */}
              <Menu theme="dark" mode="inline" defaultSelectedKeys={[currentMenuKey]}>
                {menudatas.map((value: MenuData, index: number) => (
                  <Menu.Item key={`${index+1}`}>
                    <Link to={{pathname: value.path, search: window.location.search}}>
                      <Icon type={value.icon} />
                      <span>{value.name}</span>
                    </Link>
                  </Menu.Item>
                ))}
              </Menu>
            </Sider>
            <Layout>
              {/* <Header style={{ background: '#fff', padding: 0 }} /> */}
              <Content style={{ margin: '24px 16px 0', textAlign: "left" }}>
                <div style={{ padding: 24, background: '#fff', minHeight: 360 }}>
                  <Switch>
                    {/* <Redirect exact={true} from="/" to="/app/"/> */}
                    {menudatas.map((value: MenuData, index: number) => (
                      <Route key={`route-${value.name}`} exact={true} strict={false} path={value.path} component={value.component} />
                    ))}
                  </Switch>
                </div>
              </Content>
              <Footer style={{ textAlign: 'center' }}>zilliz/mildev ©2021</Footer>
            </Layout>
          </Layout>
        </Router>
      </div>
    );
  }
}

export default App;
