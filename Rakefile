require 'fileutils'

ENV.update({
  'GOPATH' => Dir.pwd
})

DEPS = [
  'code.google.com/p/goauth2/oauth',
  'github.com/go-sql-driver/mysql',
  'github.com/googollee/go-socket.io',
  'github.com/kellegous/pork',
  'github.com/kellegous/base62',
].map do |dep|
  src = File.join('src', dep)
  file src do
    sh 'go', 'get', dep
  end
  src
end

file 'bin/four04' => DEPS + FileList['src/**'] do
  sh 'go', 'build', '-o', 'bin/four04', 'four04/fe'
end

BINS = ['bin/four04']
task :default => BINS

task :clean do
  BINS.each do |f|
    FileUtils::rm_rf(f)
  end
end