import { gql } from "@apollo/client";
import {
  Button,
  Col,
  Form,
  Input,
  Row,
  Tag
} from "antd";
import Layout, { Content } from "antd/lib/layout/layout";
import React, { useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { useAuthUserStorage } from "../../Login/UseUserStorage.hooks";
import { usePaginationQuery } from "../../PaginationHooks";
import "./Meetings.scss";

const GET_MEETINGS = gql`
  query GetMeetings($limit: Int $offset: Int, $meetingNumber: Int) {
    meetings(searchCriteria: { limit: $limit offset: $offset, number: $meetingNumber} ) {
      id
      number,
      theme,
      jointMeetingClubName,
      isAreaMeeting,
      timestamp
    }
  }
`;

const PAGE_SIZE = 15;
const DATA_KEY = 'meetings';

export default function Meetings(this: any) {
  const [meetings, loading, error, lastElementRef] = usePaginationQuery<any>(GET_MEETINGS, DATA_KEY, PAGE_SIZE)
  const history = useHistory();
  const match = useRouteMatch();
  const [authUser] = useAuthUserStorage();

  function createMeeting(id?) {
    const baseUrl = `${match.path}/create`;
    const processedUrl = (id)
      ? `${baseUrl}/${id}`
      : baseUrl;
    history.push(processedUrl)
  }

  const getSearchFilterHeader = () => {
    return (
      <Form layout="inline">
        <Form.Item name="theme">
          <Input  placeholder="Theme" />
        </Form.Item>
        { authUser?.isAdmin && 
          <Form.Item name="username">
            <Button type="primary" onClick={() => createMeeting()}>
              Create
            </Button>
          </Form.Item>
        }
      </Form>
    )
  }

  const getMeetingRow = (meeting: any) => {
    return (
      <Row className="meeting-row">
        <Col span={2}>
          <div className="meeting-number">#{meeting.number}</div>
        </Col>
        <Col span={7} className="height-100">
          <div className="meeting-label">Theme</div>
          <div className="meeting-theme">{meeting.theme}</div>
        </Col>
        <Col span={5} className="height-100">
          <div className="meeting-label">Date</div>
          <div className="meeting-value">Sat, 2nd July 2021, 09:30</div>
        </Col>
        <Col span={4} className="height-100">
          <div className="meeting-label">Location</div>
          <div className="meeting-value">{meeting.location}</div>
        </Col>
        <Col span={4} className="height-100">
          <div>
            <Tag color="cyan">cyan</Tag>
          </div>
        </Col>
        <Col span={1} className="height-100">
          { authUser?.isAdmin && 
            <div>
              <Button onClick={() => createMeeting(meeting.id)}>
                Edit
              </Button> 
            </div>
          }
        </Col>
      </Row>
    )
  }

  return (
    <Layout className="inner-layout">
      <Content>
        {getSearchFilterHeader()}
        {meetings.map((meeting: any, index) => {
          const isLastElement = meetings.length == (index + 1);
          if(isLastElement) {
            return (
              <div ref={lastElementRef} key={index}>
                {getMeetingRow(meeting)}
              </div>
            )
          }
          return (
            <div key={index}>
              {getMeetingRow(meeting)}
            </div> 
          )
        })}
        { loading &&
          <div>Loading</div>
        }
        { error && 
          <div>
            An error occurred while loading data
          </div>
        }
      </Content>
    </Layout>
  );
}
