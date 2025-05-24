<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const API_URL = 'http://localhost:8080/api'
const leagueTable = ref([])
const matches = ref([])
const currentWeek = ref(1)
const totalWeeks = ref(6)
const loading = ref(false)

const fetchLeagueTable = async () => {
  try {
    const response = await axios.get(`${API_URL}/league`)
    leagueTable.value = response.data
  } catch (error) {
    console.error('Error fetching league table:', error)
  }
}

const fetchMatches = async () => {
  try {
    const response = await axios.get(`${API_URL}/matches`)
    matches.value = response.data
  } catch (error) {
    console.error('Error fetching matches:', error)
  }
}

const simulateWeek = async () => {
  loading.value = true
  try {
    await axios.post(`${API_URL}/matches/simulate/${currentWeek.value}`)
    await Promise.all([fetchLeagueTable(), fetchMatches()])
    currentWeek.value++
  } catch (error) {
    console.error('Error simulating week:', error)
  } finally {
    loading.value = false
  }
}

const simulateAll = async () => {
  loading.value = true
  try {
    await axios.post(`${API_URL}/matches/simulate-all`)
    await Promise.all([fetchLeagueTable(), fetchMatches()])
  } catch (error) {
    console.error('Error simulating all matches:', error)
  } finally {
    loading.value = false
  }
}

const resetLeague = async () => {
  loading.value = true
  try {
    const response = await axios.post(`${API_URL}/reset`)
    leagueTable.value = response.data.stats
    matches.value = response.data.matches
    currentWeek.value = 1
  } catch (error) {
    console.error('Error resetting league:', error)
  } finally {
    loading.value = false
  }
}

const selectWeek = (week) => {
  currentWeek.value = week
}

const currentWeekMatches = computed(() => {
  return matches.value.filter(match => match.week === currentWeek.value)
})

onMounted(() => {
  fetchLeagueTable()
  fetchMatches()
})
</script>

<template>
  <div class="app">
    <header>
      <h1>Football League Simulation</h1>
    </header>
    
    <main>
      <div class="controls">
        <button @click="simulateWeek" :disabled="loading || currentWeek > totalWeeks">Simulate Next Week</button>
        <button @click="simulateAll" :disabled="loading">Simulate All Matches</button>
        <button @click="resetLeague" :disabled="loading" class="reset-button">Reset League</button>
      </div>

      <div class="content">
        <div class="league-table">
          <h2>League Table</h2>
          <table>
            <thead>
              <tr>
                <th>Pos</th>
                <th>Team</th>
                <th>P</th>
                <th>W</th>
                <th>D</th>
                <th>L</th>
                <th>GF</th>
                <th>GA</th>
                <th>GD</th>
                <th>Pts</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in leagueTable" :key="team.team_id">
                <td>{{ index + 1 }}</td>
                <td>{{ team.team_name }}</td>
                <td>{{ team.played }}</td>
                <td>{{ team.won }}</td>
                <td>{{ team.drawn }}</td>
                <td>{{ team.lost }}</td>
                <td>{{ team.goals_for }}</td>
                <td>{{ team.goals_against }}</td>
                <td>{{ team.goal_difference }}</td>
                <td>{{ team.points }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="matches">
          <h2>Matches</h2>
          <div class="week-selector">
            <button 
              v-for="week in totalWeeks" 
              :key="week"
              @click="selectWeek(week)"
              :class="{ active: currentWeek === week }"
            >
              Week {{ week }}
            </button>
          </div>
          <div class="match-list">
            <div v-for="match in currentWeekMatches" :key="match.id" class="match-card">
              <div class="match-header">Week {{ match.week }}</div>
              <div class="match-content">
                <div class="team home">
                  <span class="team-name">{{ match.home_team.name }}</span>
                  <span class="score">{{ match.home_goals }}</span>
                </div>
                <div class="vs">vs</div>
                <div class="team away">
                  <span class="team-name">{{ match.away_team.name }}</span>
                  <span class="score">{{ match.away_goals }}</span>
                </div>
              </div>
              <div class="match-status" :class="{ 'played': match.played }">
                {{ match.played ? 'Played' : 'Not Played' }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
  background-color: #f0f2f5;
  min-height: 100vh;
}

header {
  text-align: center;
  margin-bottom: 30px;
}

h1 {
  color: #1a237e;
  font-size: 2.5em;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
}

.controls {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  justify-content: center;
}

button {
  padding: 10px 20px;
  background-color: #1a237e;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1em;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

button:hover {
  background-color: #283593;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.2);
}

button:disabled {
  background-color: #9fa8da;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.league-table, .matches {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

h2 {
  color: #1a237e;
  margin-bottom: 20px;
  font-size: 1.5em;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

th, td {
  padding: 12px;
  text-align: center;
  border-bottom: 1px solid #e0e0e0;
  color: #333;
}

th {
  background-color: #e8eaf6;
  font-weight: 600;
  color: #1a237e;
}

tr:hover {
  background-color: #f5f5f5;
}

.week-selector {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.week-selector button {
  background-color: #e8eaf6;
  color: #1a237e;
}

.week-selector button.active {
  background-color: #1a237e;
  color: white;
}

.match-card {
  background: white;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
  border: 1px solid #e0e0e0;
  transition: transform 0.2s ease;
}

.match-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.match-header {
  font-size: 0.9em;
  color: #1a237e;
  margin-bottom: 10px;
  font-weight: 600;
}

.match-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.team {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.team-name {
  font-weight: 600;
  color: #333;
}

.score {
  font-size: 1.5em;
  font-weight: bold;
  color: #1a237e;
}

.vs {
  color: #666;
  font-weight: 600;
}

.match-status {
  text-align: center;
  margin-top: 10px;
  font-size: 0.9em;
  color: #666;
}

.match-status.played {
  color: #1a237e;
  font-weight: 600;
}

.reset-button {
  background-color: #d32f2f;
}

.reset-button:hover {
  background-color: #b71c1c;
}

.reset-button:disabled {
  background-color: #ef9a9a;
}

@media (max-width: 768px) {
  .content {
    grid-template-columns: 1fr;
  }
}
</style>
