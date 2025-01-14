import React, { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { Form, Row, Col, Button, Input, InputNumber, Switch } from '~/uikit'
import { Source } from '~/models/source'

const defaultFormValues = {
  source_name: '',
  host: '',
  port: '',
  password: '',
  enable_gtid: false,

  security: {
    ssl_ca_content: '',
    ssl_cert_content: '',
    ssl_key_content: '',
  },

  relay: {
    relay_enable: false,
    binlog_name: '',
    binlog_location: '',
  },
}

const CreateOrUpdateSource: React.FC<{
  onSubmit?: (values: Partial<Source>) => void
  onCancel?: () => void
  currentSource?: Source | null
}> = ({ onCancel, onSubmit, currentSource }) => {
  const [t] = useTranslation()
  const initialValues = useMemo(
    () => ({ ...defaultFormValues, ...currentSource }),
    [currentSource]
  )
  const [form] = Form.useForm()

  useEffect(() => {
    form.resetFields()
  }, [initialValues])

  return (
    <Form
      layout="vertical"
      form={form}
      onFinish={onSubmit}
      initialValues={initialValues}
    >
      <Row gutter={32}>
        <Col span={12}>
          <section>
            <h1 className="font-bold text-lg">{t('basic info')}</h1>
            <Form.Item
              name="source_name"
              label={t('source name')}
              rules={[
                { required: true, message: t('source name is required') },
              ]}
            >
              <Input
                placeholder="mysql-01"
                disabled={!!initialValues.source_name}
              />
            </Form.Item>

            <Form.Item
              name="host"
              label={t('host')}
              rules={[{ required: true, message: t('host is required') }]}
            >
              <Input placeholder="1.1.1.1" />
            </Form.Item>

            <Form.Item
              name="port"
              label={t('port')}
              rules={[{ required: true, message: t('port is required') }]}
            >
              <InputNumber placeholder="3306" />
            </Form.Item>

            <Form.Item
              name="user"
              label={t('user name')}
              rules={[{ required: true, message: t('user name is required') }]}
            >
              <Input placeholder="root" />
            </Form.Item>

            <Form.Item
              name="password"
              label={t('password')}
              rules={[{ required: true, message: t('passowrd is required') }]}
            >
              <Input type="password" />
            </Form.Item>

            <Form.Item name="enable_gtid" label="GTID" valuePropName="checked">
              <Switch defaultChecked={false} />
            </Form.Item>
          </section>
        </Col>

        <Col span={12}>
          <section>
            <h1 className="font-bold text-lg">{t('tls config (optional)')}</h1>
            <Form.Item name={['security', 'ssl_ca_content']} label="SSL CA">
              <Input.TextArea placeholder="ca.pem" />
            </Form.Item>

            <Form.Item name={['security', 'ssl_cert_content']} label="SSL CERT">
              <Input.TextArea placeholder="cert.pem" />
            </Form.Item>

            <Form.Item name={['security', 'ssl_key_content']} label="SSL KEY">
              <Input.TextArea placeholder="key.pem" />
            </Form.Item>
          </section>

          <section>
            <h1 className="font-bold text-lg">
              {t('relay config (optional)')}
            </h1>
            <Form.Item
              name={['relay', 'relay_enable']}
              label={t('enable relay')}
              valuePropName="checked"
            >
              <Switch defaultChecked={false} />
            </Form.Item>

            <Form.Item name={['relay', 'binlog_name']} label={t('binlog name')}>
              <Input placeholder={t('binlog name or gtid')} />
            </Form.Item>

            <Form.Item
              name={['relay', 'binlog_location']}
              label={t('binlog location')}
            >
              <Input placeholder="/location" />
            </Form.Item>
          </section>
        </Col>
      </Row>

      <div className="flex items-center justify-end p-4">
        <Button onClick={onCancel} className="!mr-4">
          {t('cancel')}
        </Button>
        <Button type="primary" htmlType="submit">
          {t('save')}
        </Button>
      </div>
    </Form>
  )
}

export default CreateOrUpdateSource
