import { gql, useQuery } from "@apollo/client";
import {
  Button,
  Col,
  Empty,
  Form,
  Input,
  Row, Select, Spin, Tag
} from "antd";
import Layout, { Content } from "antd/lib/layout/layout";
import classNames from "classnames";
import moment from "moment";
import React, { useEffect, useState } from "react";
import { Route, useHistory, useParams, useRouteMatch, Switch } from "react-router-dom";
import { PAGE_SEARCH_CRITERIA } from "../../Constants";
import { useAuthUserStorage } from "../../Login/UseUserStorage.hooks";
import { usePaginationQuery } from "../../PaginationHooks";
import SpeechDetail from "./SpeechDetail";
import "./Speeches.scss";

const GET_SPEECHES = gql`
  query GetSpeeches($limit: Int $offset: Int, $userId: ID, $speechTypes: [String]) {
    speeches(searchCriteria: { limit: $limit offset: $offset, userId: $userId, speechTypes: $speechTypes }) {
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

const PAGE_SIZE = 8;
const DATA_KEY = 'speeches';

export default function Speeches(this: any) {
  const [speeches, loading, error, lastElementRef, setSearchCriteria] = usePaginationQuery<any>(GET_SPEECHES, DATA_KEY, PAGE_SIZE)
  const { loading: usersLoading, error: usersError, data: usersData } = useQuery(GET_USERS, PAGE_SEARCH_CRITERIA);
  const [selectedSpeechId, setSelectedSpeechId] = useState('')
  const history = useHistory();
  const match = useRouteMatch();
  const [authUser] = useAuthUserStorage();
  const { id } = useParams() as any;
  const [searchForm] = Form.useForm();

  useEffect(() => {
    id && setSelectedSpeechId(id);
    search(searchForm.getFieldsValue())
  }, []);

  function search(searchCriteria) {
    // Skip undefined
    if(!searchCriteria.title) {
      delete searchCriteria.title;
    }
    console.log(searchCriteria);
    setSearchCriteria(searchCriteria);
  }

  function createSpeech(id?) {
    const baseUrl = `${match.path}/create`;
    const processedUrl = (id)
      ? `${baseUrl}/${id}`
      : baseUrl;
    history.push(processedUrl)
  }

  function getSearchFilterHeader(this: any) {
    const { Option } = Select;

    const speechOptions = [
      { value: 'Prepared_Speech', label: "Prepared Speech" },
      { value: 'Table_Topic', label: "Table Topic" },
      { value: 'Evaluation', label: "Evaluation" },
      { value: 'Other', label: "Other" }];

    const tagRender = (props) => {
      const { label, value, closable, onClose } = props;
      const onPreventMouseDown = event => {
        event.preventDefault();
        event.stopPropagation();
      };

      return (
        <Tag
          color={getSpeechTypeColor(value)}
          onMouseDown={onPreventMouseDown}
          closable={closable}
          onClose={onClose}
          style={{ marginRight: 3 }}
        >
          {label}
        </Tag>
      );
    }

    return (
      <Form layout="inline" onFinish={(values) => search(values)} form={searchForm}>
        <Form.Item name="userId" initialValue={authUser.id}>
          <Select loading={usersLoading} style={{ width: '200px' }}>
            { usersData?.users && !usersError &&
              usersData.users.map(user => (
                <Option value={user.id} key={user.id}>{user.firstName} {user.lastName}</Option>
              ))
            }
          </Select>
        </Form.Item>
        <Form.Item name="speechTypes" initialValue={speechOptions.map(option => option.value)}>
          <Select
            placeholder="Speech Types"
            maxTagCount={0}
            maxTagPlaceholder={(list) => `${list.length} Selected`}
            mode="multiple"
            showArrow
            tagRender={tagRender}
            style={{ width: '200px' }}
            options={speechOptions}
          />
        </Form.Item>
        <Form.Item name="title">
          <Input  placeholder="Title" />
        </Form.Item>
          <Button type="primary" htmlType="submit">
            Search
          </Button>
        { authUser?.isAdmin && 
          <Form.Item>
            <Button type="primary"  onClick={() => createSpeech()}>
              Create
            </Button>
          </Form.Item>
        }

      </Form>
    )
  }

  function getSpeechTypeColor(speechType) {
    switch(speechType) {
      case 'Prepared_Speech':
        var color = 'cyan'
        break;
      case 'Table_Topic':
        var color = 'blue';
        break;
      case 'Evaluation':
        var color = 'orange';
        break;
      default:
        var color = 'maroon';
        break; 
    }
    return color;
  }

  function getSpeechType(speechType: string) {
    return (
      <span className="tag">
        <Tag color={getSpeechTypeColor(speechType)}>{speechType.replaceAll('_', ' ')}</Tag>
      </span>
    )
  }

  function getSpeechRow(speech: any) {
    const speechRowClass = () => {
      const baseClass = `speech-row`;
      return (speech.id == selectedSpeechId)
        ? `${baseClass} selected`
        : baseClass;
    };

    const selectSpeech = (id) => {
      setSelectedSpeechId(id);

      const urlPath = (match.path === '/speeches')
        ? `${match.path}/${id}`
        : id;

      history.push(urlPath);
    }

    return (
      <Row className={speechRowClass()} onClick={() => selectSpeech(speech.id)}>
        <div className="date">{moment.unix(speech.timestamp).format("MMM DD")}</div>
        <Col span={24}>
          <div className="user">{speech.user.firstName} {speech.user.lastName}</div>
        </Col>
        <Col span={24}>
          <div className="title">{speech.title}</div>
        </Col>
        {/* <Col span={24}>
          <div>Location</div>
          <div>{speech.location}</div>
        </Col> */}
        <Col span={24}>
          {getSpeechType(speech.speechType)}
        </Col>
      </Row>
    )
  }

  return (
    <Layout className="revert-header-margin inner-layout--xl speeches-viewer">
      <Content className="speeches-content">
        <Row className="search-row">
          <Col span={24}>
            {getSearchFilterHeader()}
          </Col>
        </Row>
        <Row className="speeches-section">
          <Col span={6} className={classNames({ 'speech-list': true, 'loading-or-error-message': speeches.length == 0})}>
            <div>
              {speeches.map((speech: any, index) => {
                if(speeches.length == (index + 1)) {
                  return (
                    <div ref={lastElementRef} key={index}>
                      {getSpeechRow(speech)}
                    </div>
                  )
                }
                return (
                  <div key={index}>
                    {getSpeechRow(speech)}
                  </div> 
                )
              })}
              { loading &&
                <div className={classNames({'loading-or-error-with-data': speeches.length > 0})}>
                  <Spin size="large" delay={200}>
                  </Spin>
                </div>
              }
              { error && 
                <div className={classNames({'loading-or-error-with-data': speeches.length > 0})}>
                  An error occurred while loading data
                </div>
              }
              { !loading && !error && speeches.length == 0 && 
                <div>
                  <Empty description={<span>No Speeches For Search</span>} image={Empty.PRESENTED_IMAGE_SIMPLE} />
                </div>
              }
            </div>
          </Col>
          <Col span={18} className="speeches-section-child-wrapper">
            <Switch>
              <Route path={`${match.path}/:id`}>
                <SpeechDetail />
              </Route>
              <Route path={`${match.path}`}>
                <SpeechDetail />
              </Route>
            </Switch>
          </Col>
        </Row>
      </Content>
    </Layout>
  );
}
