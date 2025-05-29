import { useState } from 'react';

export default function TaskCard({ task, onEdit, onDelete }) {
  const [editingName, setEditingName] = useState(false);
  const [editingDescription, setEditingDescription] = useState(false);
  const [nameInput, setNameInput] = useState(task.name);
  const [descriptionInput, setDescriptionInput] = useState(task.description || '');

  const handleEditDeadline = () => {
    const deadline = prompt('Новый дедлайн (в формате 2025-06-01T12:00:00Z):', task.deadline || '');
    if (deadline !== null) {
      onEdit(task.id, { deadline });
    }
  };

  const handleNameBlur = () => {
    setEditingName(false);
    if (nameInput !== task.name) {
      onEdit(task.id, { name: nameInput });
    }
  };

  const handleDescriptionBlur = () => {
    setEditingDescription(false);
    if (descriptionInput !== task.description) {
      onEdit(task.id, { description: descriptionInput });
    }
  };

  return (
    <div style={{ border: '1px dashed gray', margin: '4px', padding: '4px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        {editingName ? (
          <input
            value={nameInput}
            onChange={(e) => setNameInput(e.target.value)}
            onBlur={handleNameBlur}
            onKeyDown={(e) => e.key === 'Enter' && handleNameBlur()}
            autoFocus
          />
        ) : (
          <strong onClick={() => setEditingName(true)} style={{ cursor: 'pointer' }}>
            {task.name}
          </strong>
        )}
        <label>
          <input
            type="checkbox"
            checked={task.done}
            onChange={() => onEdit(task.id, { done: !task.done })}
          />{' '}
          Сделано
        </label>
      </div>

      {editingDescription ? (
        <textarea
          value={descriptionInput}
          onChange={(e) => setDescriptionInput(e.target.value)}
          onBlur={handleDescriptionBlur}
          onKeyDown={(e) => e.key === 'Enter' && handleDescriptionBlur()}
          autoFocus
          style={{ width: '100%' }}
        />
      ) : (
        <div onClick={() => setEditingDescription(true)} style={{ cursor: 'pointer' }}>
          <small>{task.description || '—'}</small>
        </div>
      )}

      {task.deadline && (
        <div>
          <small>⏰ {new Date(task.deadline).toLocaleString()}</small>
        </div>
      )}

      <div style={{ marginTop: '4px' }}>
        <button onClick={handleEditDeadline}>Дедлайн</button>
        <button onClick={() => onDelete(task.id)}>Удалить</button>
      </div>
    </div>
  );
}
