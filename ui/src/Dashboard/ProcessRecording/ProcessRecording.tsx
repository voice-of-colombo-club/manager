import { Button, Checkbox, Form, Input, Select, Space, Switch, TimePicker } from 'antd'
import Layout, { Content } from 'antd/lib/layout/layout'
import React from 'react'
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { gql, useMutation, useQuery } from '@apollo/client';
import { PAGE_SEARCH_CRITERIA } from '../../Constants';
import { openNotificationWithIcon } from '../../notifications';
import moment from 'moment';

const PROCESS_RECORDING = gql`
  mutation ProcessRecording($metadata: RecordingProcessInput!) {
    processRecording(input: $metadata)
  }
`;

const GET_SPEECHES = gql`
  query GetSpeeches($limit: Int $offset: Int) {
    speeches(searchCriteria: { limit: $limit offset: $offset }) {
      id
      speechType,
      timestamp,
      title,
      user {
        id,
        firstName,
        lastName
      },
      meeting {
        id,
        theme
      }
    }
  }
`;

const T = 'HH:mm';

export default function ProcessRecording() {
  const [form] = Form.useForm();
  const { loading, error, data } = useQuery(GET_SPEECHES, PAGE_SEARCH_CRITERIA);
  const [processRecording] = useMutation(PROCESS_RECORDING);

  function getToDuration({from, to}) {
    const toFormatted = to.format('HH:mm:ss');
    const fromH = from.format('HH');
    const fromM = from.format('mm');
    const fromS = from.format('ss');
    const duration = moment.duration({hours: Number(fromH), minutes: Number(fromM), seconds: Number(fromS) })
    return moment(toFormatted, "HH:mm:ss").subtract(duration).format("HH:mm:ss");
  }

  async function onFinish(values: any) {
    try {
      values.speeches.forEach(speechEntry => {
        speechEntry.to = speechEntry.to.format('HH:mm:ss');
        speechEntry.from = speechEntry.from.format('HH:mm:ss');
        // speechEntry.galleryCrop = !!speechEntry.galleryCrop;
      });

      console.log(values);
      await processRecording({ variables: { metadata: values } });

      openNotificationWithIcon(
        "success",
        "Speech Create Success",
        "The speech has been successfully created."
      );
    } catch(e) {
      console.error(e);
      openNotificationWithIcon(
        "error",
        "Speech Create Error",
        "There was an error while creating the speech"
      );
    }
  };

  const {Option} = Select;
  return (
    <Layout className="inner-layout">
      <Content>
        <div className="content-widget">
          <Form
            form={form}
            name="nest-messages"
            onFinish={onFinish}
          >
            <Form.Item
              name="zoomMeetingId"
              rules={[{ required: true, message: 'Zoom Meeting Id Required' }]}
            >
              <Input placeholder="Zoom Meeting Id" />
            </Form.Item>
            <Form.List name="speeches">
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, fieldKey, ...restField }) => {
                    return (
                      <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                        <Form.Item
                          {...restField}
                          name={[name, 'from']}
                          fieldKey={[fieldKey, 'from']}
                          rules={[{ required: true, message: 'Missing from time' }]}
                        >
                          <TimePicker format={"HH:mm:ss"} use12Hours />
                        </Form.Item>
                        <Form.Item
                          {...restField}
                          name={[name, 'to']}
                          fieldKey={[fieldKey, 'to']}
                          rules={[{ required: true, message: 'Missing to time' }]}
                        >
                          <TimePicker format={"HH:mm:ss"} use12Hours />
                        </Form.Item>
                        <Form.Item
                          {...restField}
                          name={[name, 'speechId']}
                          fieldKey={[fieldKey, 'speechId']}
                          rules={[{ required: true, message: 'Missing speech' }]}
                        >
                         <Select loading={loading} style={{width: "200px"}}  onChange={() => {}}>
                            { data?.speeches &&
                              data.speeches.map(result => (
                                <Option value={result.id} key={result.id}>{result.title} - {result.timestamp}</Option>
                              ))
                            }
                          </Select> 
                        </Form.Item>
                        <Form.Item
                          {...restField}
                          name={[name, 'galleryCrop']}
                          fieldKey={[fieldKey, 'galleryCrop']}
                        >
                          <Switch />
                        </Form.Item>
                        <MinusCircleOutlined onClick={() => remove(name)} />
                      </Space>
                      )
                    }
                  )}
                  <Form.Item>
                    <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                      Add Speech
                    </Button>
                  </Form.Item>
                </>
              )}
            </Form.List>
            <Form.Item>
              <Button type="primary" htmlType="submit">
                Submit
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Content>
    </Layout>
  )
}
