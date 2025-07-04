// Git Repository Dashboard JavaScript

class GitDashboard {
    constructor() {
        this.data = null;
        this.filteredData = null;
        this.currentFilter = 'all';
        this.searchTerm = '';
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadData();
    }

    bindEvents() {
        document.getElementById('refreshBtn').addEventListener('click', () => this.loadData());
        document.getElementById('searchInput').addEventListener('input', (e) => {
            this.searchTerm = e.target.value.toLowerCase();
            this.filterAndRender();
        });

        // Filter buttons
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                this.currentFilter = e.target.getAttribute('data-filter');
                this.filterAndRender();
            });
        });

        // Sort select
        document.getElementById('sortSelect').addEventListener('change', () => {
            this.filterAndRender();
        });
    }

    async loadData() {
        this.showLoading(true);
        try {
            // Replace this with your actual API call
            // const response = await fetch('/api/repositories');
            // this.data = await response.json();

            // For demo purposes, using sample data
            const response = await fetch("/data");
            if (!response.ok) {
                throw new Error(`Response status: ${response.statusText}`)
            }
            const json = await response.json()
            console.log("Data", json)



            // await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate loading
            this.data = json
            this.filterAndRender();
            this.updateStats();
        } catch (error) {
            this.showError();
            console.error('Error loading data:', error);
        } finally {
            this.showLoading(false);
        }
    }

    filterAndRender() {
        if (!this.data) return;

        this.filteredData = this.filterRepositories();
        this.sortRepositories();
        this.renderRepositories();
    }

    filterRepositories() {
        const repos = [];

        Object.entries(this.data.origins).forEach(([origin, repoList]) => {
            repoList.forEach(repo => {
                const repoName = this.extractRepoName(origin);
                const matchesSearch = this.searchTerm === '' ||
                    repoName.toLowerCase().includes(this.searchTerm) ||
                    repo.folder_name.toLowerCase().includes(this.searchTerm);

                const matchesFilter = this.currentFilter === 'all' ||
                    (this.currentFilter === 'clean' && repo.clean) ||
                    (this.currentFilter === 'dirty' && !repo.clean);

                if (matchesSearch && matchesFilter) {
                    repos.push({
                        origin,
                        name: repoName,
                        ...repo
                    });
                }
            });
        });

        return repos;
    }

    sortRepositories() {
        const sortBy = document.getElementById('sortSelect').value;

        this.filteredData.sort((a, b) => {
            switch (sortBy) {
                case 'name':
                    return a.name.localeCompare(b.name);
                case 'lastCommit':
                    return new Date(b.last_commit.when) - new Date(a.last_commit.when);
                case 'status':
                    return a.clean - b.clean;
                case 'changes':
                    const aChanges = (a.untracked_files?.length || 0) + (a.changed_files?.length || 0);
                    const bChanges = (b.untracked_files?.length || 0) + (b.changed_files?.length || 0);
                    return bChanges - aChanges;
                default:
                    return 0;
            }
        });
    }

    renderRepositories() {
        const container = document.getElementById('repositoriesContainer');

        if (!this.filteredData || this.filteredData.length === 0) {
            container.innerHTML = '<div class="no-results">No repositories found matching your criteria.</div>';
            return;
        }

        container.innerHTML = this.filteredData.map(repo => this.createRepoCard(repo)).join('');
    }

    createRepoCard(repo) {
        return `
            <div class="repo-card">
                <div class="repo-header">
                    <div class="repo-title">
                        <h3>${repo.name}</h3>
                        <div class="repo-path">${repo.folder_name}</div>
                    </div>
                    <span class="repo-status ${repo.clean ? 'status-clean' : 'status-dirty'}">
                        ${repo.clean ? 'Clean' : 'Dirty'}
                    </span>
                </div>
                
                <div class="repo-info">
                    <div class="repo-branch">
                        <i class="fas fa-code-branch"></i>
                        <span>${repo.current_branch}</span>
                    </div>
                    <div class="repo-commit">
                        <strong>${repo.last_commit.author}</strong>
                        <div>${this.formatDate(repo.last_commit.when)}</div>
                        <div class="commit-message">"${repo.last_commit.message.trim()}"</div>
                    </div>
                </div>

                ${this.renderFilesSections(repo)}
            </div>
        `;
    }

    renderFilesSections(repo) {
        let html = '<div class="repo-files">';

        if (repo.untracked_files && repo.untracked_files.length > 0) {
            html += this.renderFileSection('Untracked Files', repo.untracked_files, 'untracked');
        }

        if (repo.changed_files && repo.changed_files.length > 0) {
            html += this.renderFileSection('Changed Files', repo.changed_files, 'changed');
        }

        if (repo.last_changed_files && repo.last_changed_files.length > 0) {
            html += this.renderFileSection('Recently Changed', repo.last_changed_files.slice(0, 5), 'recent');
        }

        html += '</div>';
        return html;
    }

    renderFileSection(title, files, type) {
        return `
            <div class="file-section">
                <h4>
                    <i class="fas ${this.getFileIcon(type)}"></i>
                    ${title}
                    <span class="file-count">${files.length}</span>
                </h4>
                <div class="file-list">
                    ${files.map(file => `
                        <div class="file-item">
                            <span class="file-name">${file.name}</span>
                            <span class="file-status status-${this.getFileStatusClass(file.status)}">${file.status}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    getFileIcon(type) {
        const icons = {
            'untracked': 'fa-question-circle',
            'changed': 'fa-edit',
            'recent': 'fa-clock'
        };
        return icons[type] || 'fa-file';
    }

    getFileStatusClass(status) {
        const statusMap = {
            '?': 'untracked',
            'A': 'added',
            'M': 'modified',
            'D': 'deleted',
            'R': 'renamed',
            'C': 'copied',
            ' ': 'unmodified'
        };
        return statusMap[status] || 'unknown';
    }

    updateStats() {
        if (!this.data) return;

        let totalRepos = 0;
        let dirtyRepos = 0;
        let untrackedFiles = 0;
        let changedFiles = 0;

        Object.values(this.data.origins).forEach(repoList => {
            repoList.forEach(repo => {
                totalRepos++;
                if (!repo.clean) dirtyRepos++;
                untrackedFiles += repo.untracked_files?.length || 0;
                changedFiles += repo.changed_files?.length || 0;
            });
        });

        document.getElementById('totalRepos').textContent = totalRepos;
        document.getElementById('dirtyRepos').textContent = dirtyRepos;
        document.getElementById('untrackedFiles').textContent = untrackedFiles;
        document.getElementById('changedFiles').textContent = changedFiles;
    }

    extractRepoName(origin) {
        const match = origin.match(/[\/:]([^\/]+)\.git$/);
        return match ? match[1] : origin.split('/').pop() || origin;
    }

    formatDate(dateString) {
        return new Date(dateString).toLocaleString();
    }

    showLoading(show) {
        const spinner = document.getElementById('loadingSpinner');
        const mainContent = document.querySelector('.main-content');

        if (show) {
            spinner.style.display = 'flex';
            mainContent.style.display = 'none';
        } else {
            spinner.style.display = 'none';
            mainContent.style.display = 'block';
        }
    }

    showError() {
        document.getElementById('errorMessage').style.display = 'block';
        document.querySelector('.main-content').style.display = 'none';
    }
}

// Initialize the dashboard when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new GitDashboard();
});