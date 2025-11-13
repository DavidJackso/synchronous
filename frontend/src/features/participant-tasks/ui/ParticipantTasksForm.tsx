import { useState } from 'react';
import { Input, Button, List, Typography, Card, message } from 'antd';
import { PlusOutlined, DeleteOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { tasksApi, getErrorMessage } from '@/shared/api';
import type { Task } from '@/shared/api';
import './ParticipantTasksForm.css';

const { Text, Title } = Typography;

interface ParticipantTasksFormProps {
  sessionId: string;
  initialTasks: Task[];
  onTasksChange?: (tasks: Task[]) => void;
}

/**
 * Form for participants to create and manage their own tasks in a group session
 */
export const ParticipantTasksForm = ({ 
  sessionId, 
  initialTasks,
  onTasksChange 
}: ParticipantTasksFormProps) => {
  const [tasks, setTasks] = useState<Task[]>(initialTasks);
  const [currentTask, setCurrentTask] = useState('');
  const [isAdding, setIsAdding] = useState(false);

  const handleAddTask = async () => {
    if (!currentTask.trim()) return;

    setIsAdding(true);
    try {
      const newTask = await tasksApi.addTask(sessionId, currentTask.trim());
      const updatedTasks = [...tasks, newTask.task];
      setTasks(updatedTasks);
      setCurrentTask('');
      onTasksChange?.(updatedTasks);
      message.success('Задача добавлена!');
    } catch (error) {
      console.error('[ParticipantTasksForm] Failed to add task:', error);
      message.error(`Ошибка: ${getErrorMessage(error)}`);
    } finally {
      setIsAdding(false);
    }
  };

  const handleDeleteTask = async (taskId: string) => {
    try {
      await tasksApi.deleteTask(sessionId, taskId);
      const updatedTasks = tasks.filter((t) => t.id !== taskId);
      setTasks(updatedTasks);
      onTasksChange?.(updatedTasks);
      message.success('Задача удалена');
    } catch (error) {
      console.error('[ParticipantTasksForm] Failed to delete task:', error);
      message.error(`Ошибка: ${getErrorMessage(error)}`);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTask();
    }
  };

  return (
    <Card className="participant-tasks-form">
      <Title level={5}>Ваши задачи на сессию</Title>
      <Text type="secondary" style={{ display: 'block', marginBottom: '16px' }}>
        {tasks.length === 0 
          ? 'Добавьте задачи, которые планируете выполнить'
          : `Задач: ${tasks.length}`
        }
      </Text>

      <div className="participant-tasks-form__input">
        <Input
          placeholder="Что планируете сделать?"
          value={currentTask}
          onChange={(e) => setCurrentTask(e.target.value)}
          onKeyPress={handleKeyPress}
          size="large"
          disabled={isAdding}
          suffix={
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAddTask}
              disabled={!currentTask.trim()}
              loading={isAdding}
              shape="circle"
            />
          }
        />
      </div>

      {tasks.length > 0 && (
        <List
          className="participant-tasks-form__list"
          dataSource={tasks}
          renderItem={(task) => (
            <List.Item
              className="participant-tasks-form__item"
              actions={[
                <Button
                  key="delete"
                  type="text"
                  danger
                  icon={<DeleteOutlined />}
                  onClick={() => handleDeleteTask(task.id)}
                />,
              ]}
            >
              <div className="participant-tasks-form__content">
                <CheckCircleOutlined className="participant-tasks-form__icon" />
                <Text>{task.title}</Text>
              </div>
            </List.Item>
          )}
        />
      )}
    </Card>
  );
};
