import { gql, useLazyQuery, useMutation, useQuery } from "@apollo/client";
import {
  Button,
  DatePicker,
  Form,
  Input,
  InputNumber,
} from "antd";
import Checkbox from "antd/lib/checkbox/Checkbox";
import Layout, { Content } from "antd/lib/layout/layout";
import moment from "moment";
import React, { useEffect } from "react";
import { useParams } from "react-router-dom";
import { openNotificationWithIcon } from "../notifications";
import "./CreateMeeting.scss";


const LOAD_MEETING = gql`
  query GetMeeting($id: String!) {
    meeting(id: $id) {
      id,
      theme,
      number,
      location,
      jointMeetingClubName,
      isAreaMeeting,
      timestamp
    }
  }
`;

const CREATE_MEETING = gql`
  mutation CreateMeeting($meeting: CreateMeeting!) {
    createMeeting(input: $meeting) {
      id
      number
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

export default function CreateMeeting() {
  const [createMeeting] = useMutation(CREATE_MEETING);
  const [loadMeeting, { loading: editLoading, error: editError, data: editData }] = useLazyQuery(LOAD_MEETING);
  const [form] = Form.useForm()
  const { id } = useParams() as any;

  useEffect(() => {
    if(id) {
      loadMeeting({ variables: { id: id }});
    }
  }, []);

  const setMeeting = () => {
    if(editData?.meeting) {
      const {meeting} = editData;
      const processedMeeting = {
        ...meeting,
        timestamp: moment.unix(meeting.timestamp)
      }
      form.setFieldsValue(processedMeeting);
    }
  }

  const onFinish = async (values: any) => {
    try {
      const processed = {
        ...values,
        isAreaMeeting: !!values.isAreaMeeting,
        timestamp: (values.timestamp as moment.Moment).unix(),
      };

      await createMeeting({
        variables: {
          meeting: processed,
        },
      });

      openNotificationWithIcon(
        "success",
        "Meeting Create Success",
        "The meeting has been successfully created."
      );
    } catch {
      openNotificationWithIcon(
        "error",
        "Meeting Create Error",
        "There was an error while creating the meeting"
      );
    }
  };

  setMeeting();

  const getForm = () => {
    return (<Form
            form={form}
            name="nest-messages"
            onFinish={onFinish}
            validateMessages={validateMessages}
          >
            <Form.Item
              hasFeedback
              name={"number"}
              label="Meeting Number"
              rules={[{ type: "number", min: 1, max: 10000 }]}
            >
              <InputNumber />
            </Form.Item>
            <Form.Item
              hasFeedback
              name={"theme"}
              label="Theme"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              hasFeedback
              name={"location"}
              rules={[{ required: true }]}
              label="Location"
            >
              <Input />
            </Form.Item>
            <Form.Item
              hasFeedback
              name={"timestamp"}
              rules={[{ required: true }]}
              label="Date & Time"
            >
              <DatePicker showTime />
            </Form.Item>
            <Form.Item
              name={"jointMeetingClubName"}
              label="Joint Club Name"
              hasFeedback
            >
              <Input />
            </Form.Item>
            <Form.Item name={"isAreaMeeting"}>
              <Checkbox checked>Is Area Meeting</Checkbox>
            </Form.Item>
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
              >
                Submit
              </Button>
            </Form.Item>
          </Form>)
  }

  return (
    <Layout className="inner-layout">
      <Content>
        <div className="content-widget">
          { editLoading && 
            <div>Loading</div>
          }
          { editError &&
            <div>Error</div>
          }
          { !editLoading && !editError && getForm() }
        </div>
      </Content>
    </Layout>
  );
}
