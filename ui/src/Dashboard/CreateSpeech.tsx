import { gql, useMutation, useQuery } from "@apollo/client";
import {
  Button,
  DatePicker,
  Form,
  Input,
  InputNumber,
  Select,
  notification,
} from "antd";
import Checkbox from "antd/lib/checkbox/Checkbox";
import Layout, { Content } from "antd/lib/layout/layout";
import moment from "moment";
import React, { useState } from "react";
import { useEffect } from "react";
import { PAGE_SEARCH_CRITERIA } from "../Constants";
import { openNotificationWithIcon } from "../notifications";
import "./CreateSpeech.scss";

const { Option } = Select;

const SAVE_SPEECH = gql`
  mutation SaveSpeech($speech: SaveSpeech!) {
    saveSpeech(input: $speech) {
      id
    }
  }
`;

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

const GET_USERS = gql`
  query GetUsers($limit: Int $offset: Int) {
    users(searchCriteria: { limit: $limit offset: $offset } ) {
      id
      firstName,
      lastName,
      email,
      isAdmin
    }
  }
`;

/* eslint-disable no-template-curly-in-string */
const validateMessages = {
  required: "${label} is required!",
  types: {
    email: "${label} is not a valid email!",
    number: "${label} is not a valid number!",
  },
  number: {
    range: "${label} must be between ${min} and ${max}",
  },
};
/* eslint-enable no-template-curly-in-string */

export default function CreateSpeech() {

  const [saveSpeech] = useMutation(SAVE_SPEECH);
  const { loading: meetingsLoading, error: meetingError, data: meetingsData } = useQuery(GET_MEETINGS, PAGE_SEARCH_CRITERIA);
  const { loading: usersLoading, error: usersError, data: usersData } = useQuery(GET_USERS, PAGE_SEARCH_CRITERIA);
  const [isMeetingOther, setIsMeetingOther] = useState(false);
  const [form] = Form.useForm()
  
  const onFinish = async (values: any) => {
    try {
      const processed = {
        ...values,
        timestamp: (values.timestamp as moment.Moment).unix(),
      };

      if(!processed.meetingId) {
        delete processed.meetingId
      }

      await saveSpeech({ variables: { speech: processed } });

      openNotificationWithIcon(
        "success",
        "Speech Create Success",
        "The speech has been successfully created."
      );
    } catch {
      openNotificationWithIcon(
        "error",
        "Speech Create Error",
        "There was an error while creating the speech"
      );
    }
  };

  const onMeetingChange = (meetingId) => {
    if(!meetingId) {
      setIsMeetingOther(true);
      form.setFieldsValue({ timestamp: null })
      return;
    }
    const meeting = meetingsData.meetings.find((entry) => meetingId == entry.id);
    form.setFieldsValue({ timestamp: moment.unix(meeting.timestamp) });
  }

  return (
    <Layout className="inner-layout">
      <Content>
        <div className="content-widget">
          <Form
            form={form}
            name="nest-messages"
            onFinish={onFinish}
            validateMessages={validateMessages}
          >
            <Form.Item name="meetingId" label="Meeting">
              <Select loading={meetingsLoading} onChange={onMeetingChange}>
                <Option value="" key="other">Other</Option>
                { meetingsData?.meetings &&
                  meetingsData.meetings.map(result => (
                    <Option value={result.id} key={result.id}>{result.number} - {result.theme}</Option>
                    ))
                }
              </Select>
            </Form.Item>
            <Form.Item name="userId" label="User" rules={[{ required: true }]}>
              <Select loading={usersLoading}>
                { usersData?.users &&
                  usersData.users.map(user => (
                    <Option value={user.id} key={user.id}>{user.firstName} {user.lastName}</Option>
                  ))
                }
              </Select>
            </Form.Item>
            <Form.Item
              hasFeedback
              name={"speechType"}
              label="Speech Type"
              rules={[{ required: true }]}
            >
              <Select>
                  <Option value="tabletopic">Table Topic</Option>
                  <Option value="preparedspeech">Prepared Speech</Option>
                  <Option value="evaluation">Evaluation</Option>
                  <Option value="other">Other</Option>
              </Select>
            </Form.Item>
            <Form.Item
                hasFeedback
                name="timestamp"
                rules={[{ required: true }]}
                label="Date & Time"
              > 
                <DatePicker disabled={!isMeetingOther} showTime />
            </Form.Item>
            { !!isMeetingOther &&
              <Form.Item
              name={"location"}
              label="Location"
              hasFeedback
              >
                <Input />
              </Form.Item> 
            }
            <Form.Item
              name={"project"}
              label="Project"
              hasFeedback
            >
              <Input />
            </Form.Item>

            <Form.Item
              name={"speechLink"}
              label="Speech Link"
              hasFeedback
            >
              <Input />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                // onClick={() =>
                //   addTodo({ variables: { meetings: { number: 8 } } })
                // }
              >
                Submit
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Content>
    </Layout>
  );
}
