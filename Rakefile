require 'fileutils'

ENV.update({
  'GOPATH' => Dir.pwd,
  'CGO_CFLAGS' => "-I#{Dir.pwd}/bld/rocksdb/include",
  'CGO_LDFLAGS' => "-L#{Dir.pwd}/bld/rocksdb",
})

DEPS = [
  'code.google.com/p/goauth2/oauth',
  'gopkg.in/igm/sockjs-go.v2/sockjs',
  'github.com/kellegous/pork',
  'github.com/kellegous/base62',
].map do |dep|
  src = File.join('src', dep)
  file src do
    sh 'go', 'get', dep
  end
  src
end

# rocksdb stuff
file 'pkg/darwin_amd64/github.com/DanielMorsing/rocksdb.a' do
  sh 'go', 'get', 'github.com/DanielMorsing/rocksdb'
end

file 'bld/rocksdb/Makefile' do
  FileUtils::mkdir('bld') unless File.exists?('bld')
  Process::wait spawn('git clone https://github.com/facebook/rocksdb.git',
    :chdir => 'bld')
end

file 'bld/rocksdb/librocksdb.dylib' => 'bld/rocksdb/Makefile' do
  Process::wait spawn('make shared_lib',
    :chdir => 'bld/rocksdb')
end

file 'bin/librocksdb.dylib' => 'bld/rocksdb/librocksdb.dylib' do
  FileUtils::cp('bld/rocksdb/librocksdb.dylib', 'bin/librocksdb.dylib')
end

file 'bin/fe' => DEPS + FileList['src/**/*.go'] do
  sh 'go', 'build', '-o', 'bin/fe', 'four04/fe'
end

BINS = [
  'bin/fe',
  'bin/librocksdb.dylib',
  'pkg/darwin_amd64/github.com/DanielMorsing/rocksdb.a'
]
task :default => BINS

task :clean do
  BINS.each do |f|
    FileUtils::rm_rf(f)
  end
end