import { Affix, Layout, Menu, Select } from "antd";
import Avatar from "antd/lib/avatar/avatar";
import SubMenu from "antd/lib/menu/SubMenu";
import axios from "axios";
import React, { useState } from "react";
import { useEffect } from "react";
import {
  Link, Route,
  Switch, useHistory, useRouteMatch
} from "react-router-dom";
import { useAuthUserStorage } from "../Login/UseUserStorage.hooks";
import { openNotificationWithIcon } from "../notifications";
import CreateMeeting from "./CreateMeeting";
import CreateSpeech from "./CreateSpeech";
import "./Dashboard.scss";
import Meetings from "./Meetings/Meetings";
import ProcessRecording from "./ProcessRecording/ProcessRecording";
import Speeches from "./Speech/Speeches";

export default function Dashboard() {
  const match = useRouteMatch();
  const history = useHistory();
  const [authUser, _, resetAuthUser] = useAuthUserStorage();
  const [selectedClub, setSelectedClub] = useState(authUser.clubs[0].id);

  // useEffect(() =>{
  //   setDefaultClub(authUser.clubs[0].id);
  // })

  async function logout() {
    try {
      await axios.post('/api/logout');
      resetAuthUser();  
      history.push('/login');
    }
    catch {
      openNotificationWithIcon(
        "error",
        "Unable To Logout",
        "There was an error logging you out."
      );
    }
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Affix offsetTop={0}>
        <Layout.Header className="header" style={{ background:"#ffff" }}>
          <div className="logo" />
          <Menu theme="light" className="menu-bar" mode="horizontal" defaultSelectedKeys={["1"]}>
            <Menu.Item key="1">
              <Link to=
                {`${match.path}meetings`}>Meetings
              </Link>
            </Menu.Item>
            <Menu.Item key="2">
              <Link to=
                  {`${match.path}speeches`}>Speeches
              </Link>
            </Menu.Item>
            <div style={{ position: 'absolute', top: 0, right: 0 }}>
              { authUser.clubs.length > 1 &&
                <span style={{ display: "inline-block" }}>
                  <Select defaultValue={selectedClub} style={{ width: "140px" }} onChange={(id) => setSelectedClub(id)}>
                    { authUser.clubs.map(club => (
                      <Select.Option value={club.id} key={club.id}>{club.clubName}</Select.Option>
                    ))}
                  </Select>
                </span>
              }
              <span style={{ display: "inline-block" }}>
                <SubMenu key="4" icon={<Avatar style={{ color: '#f56a00', backgroundColor: '#fde3cf' }}>{authUser?.firstName?.charAt(0)}</Avatar>}>
                    <Menu.Item onClick={logout}>Logout</Menu.Item>
                </SubMenu>
              </span>
            </div>
          </Menu>
        </Layout.Header>
      </Affix>
      <Layout>
        <Switch>
          <Route path={`${match.path}meetings/create/:id`}>
            <CreateMeeting />
          </Route>
          <Route path={`${match.path}meetings/create`}>
            <CreateMeeting />
          </Route>
          <Route path={`${match.path}speeches/create`}>
            <CreateSpeech />
          </Route>
          <Route path={`${match.path}meetings`}>
            <Meetings />
          </Route>
          <Route path={`${match.path}speeches`}>
            <Speeches />
          </Route>
          <Route path={`${match.path}process-recording`}>
            <ProcessRecording />
          </Route>
        </Switch>
      </Layout>
    </Layout>
  );
}
