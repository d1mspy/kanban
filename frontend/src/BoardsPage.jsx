import { useEffect, useState } from 'react';
import { DragDropContext, Droppable, Draggable } from '@hello-pangea/dnd';
import TaskCard from './TaskCard';

export default function BoardsPage({ token, onLogout }) {
  const [boards, setBoards] = useState([]);
  const [selectedBoard, setSelectedBoard] = useState(null);
  const [columns, setColumns] = useState([]);
  const [tasks, setTasks] = useState({});
  const [newBoardName, setNewBoardName] = useState('');
  const [newColumnName, setNewColumnName] = useState('');
  const [newTaskNames, setNewTaskNames] = useState({});

  const fetchBoards = async () => {
    const res = await fetch('/api/boards', {
      headers: { Authorization: `Bearer ${token}` }
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    if (!res.ok) return;

    const data = await res.json();
    setBoards(Array.isArray(data) ? data : []);
  };

  const fetchColumns = async (boardId) => {
    const res = await fetch(`/api/boards/${boardId}/columns`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    const data = await res.json();
    setColumns(Array.isArray(data) ? data : []);
    const tasksByColumn = {};
    for (const column of data) {
      const resTasks = await fetch(`/api/columns/${column.id}/tasks`, {
        headers: { Authorization: `Bearer ${token}` }
      });

      if (resTasks.status === 401) {
        onLogout();
        return;
      }

      tasksByColumn[column.id] = await resTasks.json();
    }
    setTasks(tasksByColumn);
  };

  const createBoard = async () => {
    if (!newBoardName.trim()) return;

    const res = await fetch('/api/boards', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name: newBoardName.trim() })
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    if (res.ok) {
      setNewBoardName('');
      fetchBoards();
    }
  };

  const deleteBoard = async (id) => {
    const res = await fetch(`/api/boards/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    if (res.ok) {
      setSelectedBoard(null);
      fetchBoards();
    }
  };

  const renameBoard = async (id, newName) => {
    const res = await fetch(`/api/boards/${id}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name: newName })
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    fetchBoards();
  };

  const createColumn = async () => {
    if (!newColumnName.trim()) return;

    const res = await fetch(`/api/boards/${selectedBoard.id}/columns`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name: newColumnName.trim() })
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    setNewColumnName('');
    fetchColumns(selectedBoard.id);
  };

  const createTask = async (columnId, name) => {
    if (!name || name.trim() === '') return;

    const res = await fetch(`/api/columns/${columnId}/tasks`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name })
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    fetchColumns(selectedBoard.id);
  };

  const editTask = async (taskId, patch) => {
    const res = await fetch(`/api/tasks/${taskId}`, {
      method: 'PATCH',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(patch)
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    fetchColumns(selectedBoard.id);
  };

  const deleteTask = async (taskId) => {
    const res = await fetch(`/api/tasks/${taskId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    });

    if (res.status === 401) {
      onLogout();
      return;
    }

    fetchColumns(selectedBoard.id);
  };

  useEffect(() => {
    fetchBoards();
  }, []);

  if (!selectedBoard) {
    return (
      <div>
        <h2>Мои доски</h2>
        <button onClick={onLogout}>Выйти</button>
        <ul>
          {boards.map(board => (
            <li key={board.id}>
              <strong>{board.name}</strong>{' '}
              <button onClick={() => {
                setSelectedBoard(board);
                setColumns([]);
                setTasks({});
                fetchColumns(board.id);
              }}>Открыть</button>
              <button onClick={() => {
                const name = prompt('Новое имя доски:', board.name);
                if (name) renameBoard(board.id, name);
              }}>Переименовать</button>
              <button onClick={() => deleteBoard(board.id)}>Удалить</button>
            </li>
          ))}
        </ul>
        <input
          value={newBoardName}
          onChange={e => setNewBoardName(e.target.value)}
          placeholder="Название доски"
          onKeyDown={e => {
            if (e.key === 'Enter') createBoard();
          }}
        />
        <button onClick={createBoard}>Создать доску</button>
      </div>
    );
  }

  return (
    <div>
      <h2>{selectedBoard.name}</h2>
      <button onClick={() => setSelectedBoard(null)}>⬅ Назад</button>
      <button onClick={onLogout}>Выйти</button>
      <div>
        <input
          value={newColumnName}
          onChange={e => setNewColumnName(e.target.value)}
          placeholder="Новая колонка"
          onKeyDown={e => {
            if (e.key === 'Enter') createColumn();
          }}
        />
        <button onClick={createColumn}>Добавить колонку</button>
      </div>
      <DragDropContext
        onDragEnd={async ({ source, destination, draggableId, type }) => {
          if (!destination || (source.droppableId === destination.droppableId && source.index === destination.index)) {
            return;
          }

          if (type === 'COLUMN') {
            const newColumns = Array.from(columns);
            const [moved] = newColumns.splice(source.index, 1);
            newColumns.splice(destination.index, 0, moved);

            setColumns(newColumns);

            await fetch(`/api/columns/${moved.id}`, {
              method: 'PATCH',
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({ position: destination.index + 1 })
            });

            fetchColumns(selectedBoard.id);
          } else {
            const fromCol = source.droppableId;
            const toCol = destination.droppableId;
            const task = tasks[fromCol].find(t => t.id === draggableId);

            const patch = fromCol === toCol
              ? { position: destination.index + 1 }
              : { column_id: toCol, position: destination.index + 1 };

            await fetch(`/api/tasks/${task.id}`, {
              method: 'PATCH',
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
              },
              body: JSON.stringify(patch)
            });

            fetchColumns(selectedBoard.id);
          }
        }}>
        <Droppable droppableId="columns" direction="horizontal" type="COLUMN">
          {(provided) => (
            <div
              ref={provided.innerRef}
              {...provided.droppableProps}
              style={{ display: 'flex', gap: '16px', marginTop: '16px' }}
            >
              {[...columns].sort((a, b) => a.position - b.position).map((column, index) => (
                <Draggable key={column.id} draggableId={column.id} index={index}>
                  {(provided) => (
                    <div
                      ref={provided.innerRef}
                      {...provided.draggableProps}
                      {...provided.dragHandleProps}
                      style={{
                        border: '1px solid black',
                        padding: '10px',
                        minWidth: '200px',
                        ...provided.draggableProps.style
                      }}
                    >
                      <h3>
                        {column.name}{' '}
                        <button onClick={() => {
                          const newName = prompt('Новое имя колонки:', column.name);
                          if (newName) {
                            fetch(`/api/columns/${column.id}`, {
                              method: 'PATCH',
                              headers: {
                                Authorization: `Bearer ${token}`,
                                'Content-Type': 'application/json'
                              },
                              body: JSON.stringify({ name: newName })
                            }).then(() => fetchColumns(selectedBoard.id));
                          }
                        }}>Переименовать</button>{' '}
                        <button onClick={() => {
                          if (confirm('Удалить колонку?')) {
                            fetch(`/api/columns/${column.id}`, {
                              method: 'DELETE',
                              headers: { Authorization: `Bearer ${token}` }
                            }).then(() => fetchColumns(selectedBoard.id));
                          }
                        }}>Удалить</button>
                      </h3>

                      <input
                        value={newTaskNames[column.id] || ''}
                        onChange={e => setNewTaskNames({ ...newTaskNames, [column.id]: e.target.value })}
                        placeholder="Новая задача"
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            createTask(column.id, newTaskNames[column.id]);
                            setNewTaskNames({ ...newTaskNames, [column.id]: '' });
                          }
                        }}
                      />
                      <button onClick={() => {
                        createTask(column.id, newTaskNames[column.id]);
                        setNewTaskNames({ ...newTaskNames, [column.id]: '' });
                      }}>+</button>

                      <Droppable droppableId={column.id} type="TASK">
                        {(provided) => (
                          <div ref={provided.innerRef} {...provided.droppableProps}>
                            {(tasks[column.id] || []).sort((a, b) => a.position - b.position).map((task, i) => (
                              <Draggable key={task.id} draggableId={task.id} index={i}>
                                {(provided) => (
                                  <div ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                                    <TaskCard task={task} onEdit={editTask} onDelete={deleteTask} />
                                  </div>
                                )}
                              </Draggable>
                            ))}
                            {provided.placeholder}
                          </div>
                        )}
                      </Droppable>
                    </div>
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </div>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  );
}
