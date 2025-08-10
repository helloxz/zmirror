<template>
  <div class="home-container">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <div class="hero-text">
          <h1 class="hero-title">
            <span class="gradient-text">ZMirror</span>
          </h1>
          <h2 class="hero-subtitle">Docker镜像代理管理系统</h2>
          <p class="hero-description">
            一个简单易用的Docker镜像代理管理系统，支持多镜像源配置、白名单管理和访问控制。
            帮助企业和个人更好地管理Docker镜像拉取，提高开发效率。
          </p>
          
          <div class="hero-actions">
            <el-button 
              type="primary" 
              size="large" 
              @click="handleMainAction"
              :icon="authStore.isLoggedIn ? Monitor : Lock"
            >
              {{ authStore.isLoggedIn ? '管理后台' : '立即登录' }}
            </el-button>
            
            <el-button 
              size="large" 
              @click="openGithub"
              
            >
              <el-icon style="margin-right: 6px;"><Star /></el-icon>
              GitHub
            </el-button>
          </div>
        </div>
        
        <div class="hero-image">
          <div class="docker-icon">
            <svg t="1754723796933" class="icon" viewBox="0 0 1423 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="6390" width="200" height="200"><path d="M285.924091 824.783652c-40.580953 0-77.378948-33.238872-77.378948-73.54594s33.238872-73.730847 77.423437-73.730847c44.362522 0 77.734861 33.238872 77.734861 73.684968s-36.886973 73.54594-77.77935 73.641869z m949.952168-401.222992c-8.009417-58.868728-44.495989-106.791763-92.55249-143.589758l-18.688176-14.81763-15.084565 18.376753c-29.322446 33.238872-40.936866 92.152089-37.376352 136.158698 3.559123 33.32785 14.239272 66.433255 32.88296 92.196579-15.084564 7.697994-33.683762 14.81763-48.056502 22.383547-33.817229 11.079161-66.700189 14.81763-99.672127 14.81763H5.686255l-3.560513 21.934485c-7.119636 70.125845 3.560513 143.588368 33.325069 210.023013l14.458937 25.807813v3.559123c88.991977 147.326837 247.399365 213.582136 419.911169 213.582136 331.89795 0 604.077709-143.589758 733.120942-452.838618 84.543074 3.69259 169.931439-18.376753 210.023013-99.449681l10.68293-18.376753-17.795615-11.079161c-48.056502-29.322446-113.910009-33.238872-169.086147-18.376753z m-475.087306-58.868728h-144.033258v143.589758h144.172286V364.602954z m0-180.521221h-144.033258v143.589759h144.172286V184.348668z m0-184.170711h-144.033258v143.589758h144.172286V0.004171z m176.204393 364.691932H793.7595v143.589758h143.589758V364.602954z m-533.868447 0h-143.143478v143.589758h143.678736V364.602954z m179.763516 0H440.504032v143.589758h143.232455V364.602954z m-357.750252 0H82.976225v143.589758h144.033259V364.602954zM582.891195 184.170711H440.504032v143.589759h143.232455V184.348668z m-180.654687 0H260.203866v143.589759H403.348734V184.348668z" fill="#1488C6" p-id="6391"></path></svg>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <div class="container">
        <h2 class="section-title">核心功能</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">
              <el-icon size="48" color="#409EFF"><Connection /></el-icon>
            </div>
            <h3>多镜像源管理</h3>
            <p>支持配置多个Docker镜像源，实现负载均衡和故障切换</p>
          </div>
          
          <div class="feature-card">
            <div class="feature-icon">
              <el-icon size="48" color="#67C23A"><Key /></el-icon>
            </div>
            <h3>白名单控制</h3>
            <p>精确控制访问权限，支持IP白名单和用户认证</p>
          </div>
          
          <div class="feature-card">
            <div class="feature-icon">
              <el-icon size="48" color="#E6A23C"><Document /></el-icon>
            </div>
            <h3>访问日志</h3>
            <p>详细的访问日志记录，便于监控和审计</p>
          </div>
          
          <div class="feature-card">
            <div class="feature-icon">
              <el-icon size="48" color="#F56C6C"><UserFilled /></el-icon>
            </div>
            <h3>用户管理</h3>
            <p>多用户支持，灵活的权限管理体系</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Usage Section -->
    <section class="usage-section">
      <div class="container">
        <h2 class="section-title">快速使用</h2>
        <div class="usage-content">
          <div class="usage-steps">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <h3>配置镜像源</h3>
                <p>在管理后台添加您的Docker镜像源地址</p>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <h3>设置白名单</h3>
                <p>配置允许访问的IP地址或用户</p>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <h3>使用代理</h3>
                <p>将Docker配置指向本系统即可开始使用</p>
              </div>
            </div>
          </div>
          
          <div class="code-example">
            <h3>配置示例</h3>
            <pre><code># 在Docker daemon配置中添加
{
  "registry-mirrors": [
    "http://your-server:8080"
  ]
}</code></pre>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
      <div class="container">
        <p>&copy; 2024 ZMirror. 开源项目，欢迎贡献代码。</p>
        <div class="footer-links">
          <el-button text @click="openGithub">
            <el-icon><Star /></el-icon>
            GitHub
          </el-button>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { 
  Monitor,
  Lock,
  Star,
  Connection,
  Key,
  Document,
  UserFilled
} from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

const handleMainAction = () => {
  if (authStore.isLoggedIn) {
    router.push('/dashboard')
  } else {
    router.push('/login')
  }
}

const openGithub = () => {
  window.open('https://github.com/helloxz/zmirror', '_blank')
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* Hero Section */
.hero-section {
  min-height: 100vh;
  display: flex;
  align-items: center;
  padding: 0 24px;
  position: relative;
  overflow: hidden;
}

.hero-content {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 80px;
  align-items: center;
  width: 100%;
}

.hero-text {
  color: white;
}

.hero-title {
  font-size: 4rem;
  font-weight: 800;
  margin: 0 0 16px 0;
  line-height: 1.1;
}

.gradient-text {
  background: linear-gradient(45deg, #fff, #a8edea);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 24px 0;
  opacity: 0.9;
}

.hero-description {
  font-size: 1.125rem;
  line-height: 1.7;
  margin: 0 0 40px 0;
  opacity: 0.8;
  max-width: 500px;
}

.hero-actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.hero-image {
  display: flex;
  justify-content: center;
  align-items: center;
}

.docker-icon {
  animation: float 3s ease-in-out infinite;
  filter: drop-shadow(0 10px 20px rgba(0,0,0,0.2));
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

/* Features Section */
.features-section {
  padding: 100px 24px;
  background: white;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 700;
  text-align: center;
  margin: 0 0 64px 0;
  color: #2c3e50;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 32px;
}

.feature-card {
  background: white;
  padding: 40px 24px;
  border-radius: 16px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 40px rgba(0,0,0,0.15);
}

.feature-icon {
  margin-bottom: 24px;
}

.feature-card h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 16px 0;
  color: #2c3e50;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
  margin: 0;
}

/* Usage Section */
.usage-section {
  padding: 100px 24px;
  background: #f8f9fa;
}

.usage-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 64px;
  align-items: start;
}

.usage-steps {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.step {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.step-number {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 1.2rem;
  flex-shrink: 0;
}

.step-content h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #2c3e50;
}

.step-content p {
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.code-example {
  background: #2c3e50;
  border-radius: 12px;
  padding: 32px;
  color: white;
}

.code-example h3 {
  margin: 0 0 20px 0;
  font-size: 1.25rem;
}

.code-example pre {
  margin: 0;
  overflow-x: auto;
}

.code-example code {
  font-family: 'Fira Code', monospace;
  font-size: 0.9rem;
  line-height: 1.6;
}

/* Footer */
.footer {
  background: #2c3e50;
  color: white;
  padding: 40px 24px;
  text-align: center;
}

.footer .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer p {
  margin: 0;
  opacity: 0.8;
}

.footer-links {
  display: flex;
  gap: 16px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .hero-content {
    grid-template-columns: 1fr;
    gap: 40px;
    text-align: center;
  }
  
  .usage-content {
    grid-template-columns: 1fr;
    gap: 40px;
  }
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-subtitle {
    font-size: 1.25rem;
  }
  
  .section-title {
    font-size: 2rem;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .footer .container {
    flex-direction: column;
    gap: 16px;
  }
  
  .hero-actions {
    justify-content: center;
  }
  
  .features-section,
  .usage-section {
    padding: 60px 16px;
  }
  
  .hero-section {
    padding: 0 16px;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: 2rem;
  }
  
  .hero-subtitle {
    font-size: 1.1rem;
  }
  
  .hero-description {
    font-size: 1rem;
  }
  
  .feature-card {
    padding: 24px 16px;
  }
  
  .code-example {
    padding: 20px;
  }
}
</style>
