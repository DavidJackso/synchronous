import { Avatar, Typography, Button } from 'antd';
import { ArrowLeftOutlined, UserOutlined } from '@ant-design/icons';
import './Header.css';

const { Text, Title } = Typography;

interface HeaderProps {
  variant?: 'home' | 'page';
  userName?: string;
  pageTitle?: string;
  onBack?: () => void;
  avatarUrl?: string;
}

export const Header = ({
  variant = 'home',
  userName = 'Пользователь',
  pageTitle,
  onBack,
  avatarUrl,
}: HeaderProps) => {
  return (
    <header className="header">
      {variant === 'home' ? (
        <>
          <div className="header__left">
            <Text type="secondary" className="header__greeting">
              Привет!
            </Text>
            <Title level={4} className="header__username">
              {userName}
            </Title>
          </div>
          <div className="header__right">
            <Avatar size={44} src={avatarUrl} icon={<UserOutlined />} />
          </div>
        </>
      ) : (
        <>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={onBack}
            className="header__back"
          />
          <Title level={4} className="header__page-title">
            {pageTitle}
          </Title>
          <div className="header__right">
            <Avatar size={36} src={avatarUrl} icon={<UserOutlined />} />
          </div>
        </>
      )}
    </header>
  );
};
