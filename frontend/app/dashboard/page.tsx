'use client';

import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useState, FormEvent } from 'react';
import { fetchApi } from '@/lib/api';

interface Project {
  id: string;
  userID: string;
  name: string;
  description?: string | null;
  createdAt: string;
  updatedAt: string;
}

export default function DashboardPage() {
  const { token, logout, isLoading } = useAuth();
  const router = useRouter();
  const [projects, setProjects] = useState<Project[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loadingProjects, setLoadingProjects] = useState(true);

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [newProjectDescription, setNewProjectDescription] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [createFormError, setCreateFormError] = useState<string | null>(null);

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [editProjectName, setEditProjectName] = useState('');
  const [editProjectDescription, setEditProjectDescription] = useState('');
  const [isUpdating, setIsUpdating] = useState(false);
  const [editFormError, setEditFormError] = useState<string | null>(null);

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [deletingProject, setDeletingProject] = useState<Project | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const loadProjects = async () => {
    if (!token) return;
    setLoadingProjects(true);
    setError(null);
    try {
      const data = await fetchApi<Project[]>('/projects', token);
      setProjects(data || []);
    } catch (err: unknown) {
      console.error('Failed to fetch projects:', err);
      setError(err instanceof Error ? err.message : 'Nie udało się załadować projektów.');
    } finally {
      setLoadingProjects(false);
    }
  };


  useEffect(() => {
    if (!isLoading && !token) {
      router.push('/login');
    }
    if (token) {
      loadProjects();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, isLoading, router]); // Usunięto loadProjects z zależności, by uniknąć pętli

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  // --- NOWE FUNKCJE OBSŁUGI MODALA I FORMULARZA ---
  const handleAddProjectClick = () => {
    setNewProjectName('');
    setNewProjectDescription('');
    setCreateFormError(null);
    setIsCreateModalOpen(true);
  };

  const handleCreateProject = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token) return;

    setIsCreating(true);
    setCreateFormError(null);

    try {
      const newProjectData = {
        name: newProjectName,
        description: newProjectDescription || null,
      };

      await fetchApi<Project>('/projects', token, {
        method: 'POST',
        body: JSON.stringify(newProjectData),
      });

      setIsCreateModalOpen(false);
      await loadProjects();

    } catch (err: unknown) {
        console.error('Failed to create project:', err);
        setCreateFormError(err instanceof Error ? err.message : 'Nie udało się utworzyć projektu.');
    } finally {
        setIsCreating(false);
    }
  };

  const handleEditClick = (project: Project) => {
    setEditingProject(project);
    setEditProjectName(project.name);
    setEditProjectDescription(project.description || '');
    setEditFormError(null);
    setIsEditModalOpen(true);
  };

  const handleUpdateProject = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!token || !editingProject) return;
    setIsUpdating(true);
    setEditFormError(null);
    try {
        await fetchApi<Project>(`/projects/${editingProject.id}`, token, {
            method: 'PUT',
            body: JSON.stringify({ name: editProjectName, description: editProjectDescription || null }),
        });
        setIsEditModalOpen(false);
        setEditingProject(null);
        await loadProjects();
    } catch (err: unknown) {
        console.error('Failed to update project:', err);
        setEditFormError(err instanceof Error ? err.message : 'Nie udało się zaktualizować projektu.');
    } finally {
        setIsUpdating(false);
    }
  };

    const handleDeleteClick = (project: Project) => {
        setDeletingProject(project);
        setDeleteError(null);
        setIsDeleteModalOpen(true);
    };

  const handleConfirmDelete = async () => {
    if (!token || !deletingProject) return;
    setIsDeleting(true);
    setDeleteError(null);
    try {
        await fetchApi<null>(`/projects/${deletingProject.id}`, token, {
            method: 'DELETE',
        });
        setIsDeleteModalOpen(false);
        setDeletingProject(null);
        await loadProjects();
    } catch (err: unknown) {
        console.error('Failed to delete project:', err);
        setDeleteError(err instanceof Error ? err.message : 'Nie udało się usunąć projektu.');
         // Nie zamykaj modala przy błędzie, aby użytkownik widział komunikat
    } finally {
        setIsDeleting(false);
    }
  };

  if (isLoading) {
    return <div>Sprawdzanie sesji...</div>;
  }
   if (!token) {
      return null;
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      {/* ... (Header bez zmian) ... */}
      <header className="mb-8 flex items-center justify-between rounded-lg bg-white p-4 shadow">
        <h1 className="text-xl font-semibold text-gray-800">Panel raitometer</h1>
        <button
          onClick={handleLogout}
          className="rounded-md bg-red-600 px-3 py-1 text-sm font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
        >
          Wyloguj się
        </button>
      </header>


      <div className="rounded-lg bg-white p-6 shadow">
        <h2 className="mb-4 text-lg font-medium text-gray-700">Twoje Projekty</h2>

        {error && (
            <div className="rounded-md bg-red-50 p-4 mb-4">
              <p className="text-sm text-red-700">{error}</p>
            </div>
        )}
        {loadingProjects && <p className="text-gray-500">Ładowanie projektów...</p>}
        {!loadingProjects && projects.length === 0 && (
          <p className="text-gray-500">Nie masz jeszcze żadnych projektów.</p>
        )}
        {!loadingProjects && projects.length > 0 && (
          <ul className="space-y-3">
            {projects.map((project) => (
              <li key={project.id} className="rounded border border-gray-200 p-4">
                <div>
                    <h3 className="font-semibold text-gray-800">{project.name}</h3>
                    {project.description && (
                    <p className="mt-1 text-sm text-gray-600">{project.description}</p>
                    )}
                    <p className="mt-2 text-xs text-gray-400">
                    Utworzono: {new Date(project.createdAt).toLocaleDateString('pl-PL')}
                    </p>
                </div>
                <div className="space-x-2 flex-shrink-0">
                     <button
                        onClick={() => router.push(`/dashboard/project/${project.id}`)} // TODO: Placeholder link
                        className="rounded bg-blue-100 px-2 py-1 text-xs font-medium text-blue-700 hover:bg-blue-200"
                    >
                        Otwórz
                    </button>
                    <button
                        onClick={() => handleEditClick(project)}
                        className="rounded bg-yellow-100 px-2 py-1 text-xs font-medium text-yellow-700 hover:bg-yellow-200"
                    >
                        Edytuj
                    </button>
                     <button
                        onClick={() => handleDeleteClick(project)}
                        className="rounded bg-red-100 px-2 py-1 text-xs font-medium text-red-700 hover:bg-red-200"
                    >
                        Usuń
                    </button>
                </div>
              </li>
            ))}
          </ul>
        )}

        {/* ZMIANA TUTAJ: Przycisk otwiera modal */}
        <button
            onClick={handleAddProjectClick}
            className="mt-6 rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            + Dodaj Nowy Projekt
        </button>
      </div>

      {/* --- NOWY MODAL --- */}
      {isCreateModalOpen && (
        <div className="fixed inset-0 z-10 overflow-y-auto bg-gray-500 bg-opacity-75 transition-opacity">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <div className="relative w-full max-w-lg transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all">
              <form onSubmit={handleCreateProject}>
                <div className="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    Nowy Projekt
                  </h3>
                  <div className="mt-4 space-y-4">
                    <div>
                      <label htmlFor="projectName" className="block text-sm font-medium text-gray-700">
                        Nazwa projektu *
                      </label>
                      <input
                        type="text"
                        name="projectName"
                        id="projectName"
                        required
                        value={newProjectName}
                        onChange={(e) => setNewProjectName(e.target.value)}
                        className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>
                    <div>
                      <label htmlFor="projectDescription" className="block text-sm font-medium text-gray-700">
                        Opis (opcjonalnie)
                      </label>
                      <textarea
                        name="projectDescription"
                        id="projectDescription"
                        rows={3}
                        value={newProjectDescription}
                        onChange={(e) => setNewProjectDescription(e.target.value)}
                        className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>
                    {createFormError && (
                      <p className="text-sm text-red-600">{createFormError}</p>
                    )}
                  </div>
                </div>
                <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                  <button
                    type="submit"
                    disabled={isCreating}
                    className={`inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm ${isCreating ? 'cursor-not-allowed opacity-50' : ''}`}
                  >
                    {isCreating ? 'Tworzenie...' : 'Utwórz Projekt'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setIsCreateModalOpen(false)}
                    className="mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                  >
                    Anuluj
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {isEditModalOpen && editingProject && (
        <div className="fixed inset-0 z-10 overflow-y-auto bg-gray-500 bg-opacity-75 transition-opacity">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <div className="relative w-full max-w-lg transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all">
              <form onSubmit={handleUpdateProject}>
                <div className="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                  <h3 className="text-lg font-medium leading-6 text-gray-900">
                    Edytuj Projekt: {editingProject.name}
                  </h3>
                  <div className="mt-4 space-y-4">
                    <div>
                      <label htmlFor="editProjectName" className="block text-sm font-medium text-gray-700">
                        Nazwa projektu *
                      </label>
                      <input
                        type="text"
                        name="editProjectName"
                        id="editProjectName"
                        required
                        value={editProjectName}
                        onChange={(e) => setEditProjectName(e.target.value)}
                        className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>
                    <div>
                      <label htmlFor="editProjectDescription" className="block text-sm font-medium text-gray-700">
                        Opis (opcjonalnie)
                      </label>
                      <textarea
                        name="editProjectDescription"
                        id="editProjectDescription"
                        rows={3}
                        value={editProjectDescription}
                        onChange={(e) => setEditProjectDescription(e.target.value)}
                        className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                      />
                    </div>
                    {editFormError && (
                      <p className="text-sm text-red-600">{editFormError}</p>
                    )}
                  </div>
                </div>
                <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                  <button
                    type="submit"
                    disabled={isUpdating}
                    className={`inline-flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm ${isUpdating ? 'cursor-not-allowed opacity-50' : ''}`}
                  >
                    {isUpdating ? 'Zapisywanie...' : 'Zapisz Zmiany'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setIsEditModalOpen(false)}
                    className="mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                  >
                    Anuluj
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
      {/* --- KONIEC MODALA EDYCJI --- */}

      {/* --- NOWY MODAL POTWIERDZENIA USUNIĘCIA --- */}
       {isDeleteModalOpen && deletingProject && (
        <div className="fixed inset-0 z-10 overflow-y-auto bg-gray-500 bg-opacity-75 transition-opacity">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <div className="relative w-full max-w-md transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all">
               <div className="bg-white px-4 pb-4 pt-5 sm:p-6 sm:pb-4">
                  <div className="sm:flex sm:items-start">
                    <div className="mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                      {/* Ikona ostrzeżenia (opcjonalnie) */}
                      <svg className="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" /></svg>
                    </div>
                    <div className="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                      <h3 className="text-lg font-medium leading-6 text-gray-900">
                        Usuń Projekt
                      </h3>
                      <div className="mt-2">
                        <p className="text-sm text-gray-500">
                          Czy na pewno chcesz usunąć projekt {deletingProject.name}? Tej operacji nie można cofnąć.
                        </p>
                      </div>
                       {deleteError && (
                          <p className="mt-2 text-sm text-red-600">{deleteError}</p>
                        )}
                    </div>
                  </div>
                </div>
               <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6">
                  <button
                    type="button"
                    disabled={isDeleting}
                    onClick={handleConfirmDelete}
                    className={`inline-flex w-full justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-base font-medium text-white shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm ${isDeleting ? 'cursor-not-allowed opacity-50' : ''}`}
                  >
                    {isDeleting ? 'Usuwanie...' : 'Tak, usuń'}
                  </button>
                  <button
                    type="button"
                    onClick={() => setIsDeleteModalOpen(false)}
                    className="mt-3 inline-flex w-full justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-base font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
                  >
                    Anuluj
                  </button>
                </div>
            </div>
          </div>
        </div>
      )}
      {/* --- KONIEC MODALA POTWIERDZENIA USUNIĘCIA --- */}
    </div>
  );
}